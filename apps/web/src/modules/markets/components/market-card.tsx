import Link from "next/link"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { getMarketCategoryLabel } from "@/modules/categories"
import type { Market } from "@/lib/api"

function formatDate(value: string) {
  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date)
}

function statusColor(status: string) {
  switch (status.toLowerCase()) {
    case "open":
      return "border-green-500/30 bg-green-500/10 text-green-300"
    case "closed":
      return "border-yellow-500/30 bg-yellow-500/10 text-yellow-300"
    case "resolved":
      return "border-indigo-500/30 bg-indigo-500/10 text-indigo-300"
    case "cancelled":
      return "border-red-500/30 bg-red-500/10 text-red-300"
    default:
      return ""
  }
}

function deploymentColor(status: string) {
  switch (status) {
    case "DEPLOYED":
      return "border-green-500/20 bg-green-500/5 text-green-400"
    case "FAILED":
      return "border-red-500/20 bg-red-500/5 text-red-400"
    default:
      return "border-border bg-muted/30 text-muted-foreground"
  }
}

function deploymentLabel(status: string) {
  switch (status) {
    case "DEPLOYED":
      return "Onchain"
    case "FAILED":
      return "Deploy Failed"
    default:
      return "Not Deployed"
  }
}

export function MarketCard({ market }: { market: Market }) {
  return (
    <Card className="group transition-colors hover:border-indigo-500/30">
      <div className="px-6 pt-6">
        {market.cover_image_url ? (
          // Plain img is intentional for v1 user-provided remote URLs.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={market.cover_image_url}
            alt={market.title}
            className="h-48 w-full rounded-xl object-cover"
            loading="lazy"
          />
        ) : (
          <div className="h-48 w-full rounded-xl bg-muted" aria-hidden="true" />
        )}
      </div>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div className="min-w-0 space-y-2.5">
            <CardTitle className="text-base font-semibold leading-snug">
              <Link
                className="transition-colors hover:text-indigo-300"
                href={`/markets/${market.id}`}
              >
                {market.title}
              </Link>
            </CardTitle>

            <CardDescription className="flex flex-wrap items-center gap-1.5">
              <Badge variant="outline" className={statusColor(market.status)}>
                {market.status}
              </Badge>

              <Badge
                variant="outline"
                className="border-muted-foreground/20 bg-muted/40 text-muted-foreground"
              >
                {getMarketCategoryLabel(market.category)}
              </Badge>

              <Badge
                variant="outline"
                className={`text-xs ${deploymentColor(market.onchain_deployment_status)}`}
              >
                {deploymentLabel(market.onchain_deployment_status)}
              </Badge>
            </CardDescription>
          </div>

          <Button
            asChild
            size="sm"
            variant="outline"
            className="shrink-0 opacity-0 transition-opacity group-hover:opacity-100"
          >
            <Link href={`/markets/${market.id}`}>View</Link>
          </Button>
        </div>
      </CardHeader>

      <CardContent>
        <dl className="grid gap-4 text-sm text-muted-foreground sm:grid-cols-3">
          <div className="space-y-0.5">
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
              Collateral
            </dt>
            <dd className="font-medium text-foreground">
              {market.collateral_asset}
            </dd>
          </div>
          <div className="space-y-0.5">
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
              Chain
            </dt>
            <dd className="font-medium text-foreground">{market.chain}</dd>
          </div>
          <div className="space-y-0.5">
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
              Closes
            </dt>
            <dd className="font-medium text-foreground">
              {formatDate(market.closes_at)}
            </dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  )
}
