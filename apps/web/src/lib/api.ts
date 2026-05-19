type RequestMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE"

export type ApiResponse<TData> = {
  data: TData
  requestId: string | null
}

export type ApiErrorBody = {
  code: string
  message: string
}

type BackendErrorResponse = {
  error?: Partial<ApiErrorBody>
}

export class ApiError extends Error {
  status: number
  code: string
  requestId: string | null
  body: unknown

  constructor({
    status,
    code,
    message,
    requestId,
    body,
  }: {
    status: number
    code: string
    message: string
    requestId: string | null
    body: unknown
  }) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.code = code
    this.requestId = requestId
    this.body = body
  }
}

export type HealthResponse = {
  status: string
}

export type SchemaValidationResponse = {
  status: string
  migration_version: number
  dirty: boolean
  missing_tables: string[]
}

export type Market = {
  id: string
  creator_user_id: string
  title: string
  description: string | null
  category: string | null
  status: string
  outcome_yes_label: string
  outcome_no_label: string
  collateral_asset: string
  chain: string
  resolution_source: string | null
  opens_at: string | null
  closes_at: string
  resolved_at: string | null
  settled_at: string | null
  winning_outcome: string | null
  created_at: string
  updated_at: string
}

export type AgentMarket = {
  id: string
  title: string
  status: string
  category: string | null
  collateral_asset: string
  chain: string
  closes_at: string
  resolution_source: string | null
}

export type MarketsResponse = {
  markets: Market[]
}

export type MarketResponse = {
  market: Market
}

export type CreateMarketRequest = {
  creator_user_id: string
  title: string
  description?: string
  category?: string
  outcome_yes_label?: string
  outcome_no_label?: string
  collateral_asset?: string
  chain: string
  resolution_source?: string
  opens_at?: string
  closes_at: string
}

export type AgentMarketsResponse = {
  markets: AgentMarket[]
}

type ApiRequestOptions = Omit<RequestInit, "body" | "method"> & {
  body?: unknown
  method?: RequestMethod
}

function getApiBaseUrl() {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL?.trim()

  if (!baseUrl) {
    throw new Error("NEXT_PUBLIC_API_BASE_URL is required to call the SignalArc API")
  }

  return baseUrl.replace(/\/+$/, "")
}

function buildUrl(path: string) {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`
  return `${getApiBaseUrl()}${normalizedPath}`
}

function getRequestId(response: Response) {
  return response.headers.get("X-Request-ID")
}

async function readJson(response: Response): Promise<unknown> {
  if (response.status === 204) {
    return null
  }

  const text = await response.text()
  if (!text) {
    return null
  }

  try {
    return JSON.parse(text)
  } catch {
    return text
  }
}

function normalizeApiError(response: Response, requestId: string | null, body: unknown) {
  const backendError = body as BackendErrorResponse
  const code = backendError?.error?.code ?? `http_${response.status}`
  const message = backendError?.error?.message ?? (response.statusText || "API request failed")

  return new ApiError({
    status: response.status,
    code,
    message,
    requestId,
    body,
  })
}

export async function apiRequest<TData>(
  path: string,
  options: ApiRequestOptions = {},
): Promise<ApiResponse<TData>> {
  const headers = new Headers(options.headers)
  const hasBody = options.body !== undefined

  if (hasBody && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json")
  }

  const response = await fetch(buildUrl(path), {
    ...options,
    method: options.method ?? "GET",
    headers,
    body: hasBody ? JSON.stringify(options.body) : undefined,
  })

  const requestId = getRequestId(response)
  const body = await readJson(response)

  if (!response.ok) {
    throw normalizeApiError(response, requestId, body)
  }

  return {
    data: body as TData,
    requestId,
  }
}

export function getHealth() {
  return apiRequest<HealthResponse>("/health")
}

export function getReadiness() {
  return apiRequest<HealthResponse>("/readyz")
}

export function validateSchema() {
  return apiRequest<SchemaValidationResponse>("/schema/validate")
}

export function getMarkets() {
  return apiRequest<MarketsResponse>("/markets")
}

export function getMarket(id: string) {
  return apiRequest<MarketResponse>(`/markets/${encodeURIComponent(id)}`)
}

export function createMarket(input: CreateMarketRequest) {
  return apiRequest<MarketResponse>("/markets", {
    method: "POST",
    body: input,
  })
}

export function getAgentMarkets() {
  return apiRequest<AgentMarketsResponse>("/agent/markets")
}
