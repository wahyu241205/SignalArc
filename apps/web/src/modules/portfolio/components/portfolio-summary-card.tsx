import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { InlineErrorState } from "@/components/shared"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import type { MarketsState } from "../types"
import { truncatePortfolioId } from "../format"

export function PortfolioSummaryCard({ state }: { state: MarketsState }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Onchain Markets</CardTitle>
        <CardDescription>
          Arc Testnet deployment status for listed markets.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {state.status === "loading" ? (
          <div className="h-4 w-1/2 animate-pulse rounded bg-muted" />
        ) : null}

        {state.status === "error" ? (
          <InlineErrorState
            title="Unable to load onchain market status"
            message={state.message}
            requestId={state.requestId}
          />
        ) : null}

        {state.status === "loaded" ? (
          <div className="overflow-x-auto rounded-lg border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Market</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Contract</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {state.markets.map((market) => (
                  <TableRow key={market.id}>
                    <TableCell className="max-w-[260px] truncate" title={market.title}>
                      {market.title}
                    </TableCell>
                    <TableCell>{market.onchain_deployment_status}</TableCell>
                    <TableCell className="font-mono text-xs">
                      {market.market_contract_address
                        ? truncatePortfolioId(market.market_contract_address)
                        : "Not deployed"}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        ) : null}
      </CardContent>
    </Card>
  )
}
