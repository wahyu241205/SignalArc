import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import type { Position, Settlement } from "@/lib/api"

import { formatPortfolioDate, truncatePortfolioId } from "../format"

function PositionsTable({ positions }: { positions: Position[] }) {
  if (positions.length === 0) {
    return <p className="text-sm text-muted-foreground">No open positions.</p>
  }

  return (
    <>
      <div className="grid gap-3 md:hidden">
        {positions.map((position) => (
          <div key={position.id} className="rounded-lg border bg-card/50 p-3">
            <div className="flex items-start justify-between gap-3">
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Market
                </p>
                <p className="mt-1 font-mono text-xs" title={position.market_id}>
                  {truncatePortfolioId(position.market_id)}
                </p>
              </div>
              <div className="text-right">
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Outcome
                </p>
                <p className="mt-1 text-sm font-medium">{position.outcome}</p>
              </div>
            </div>
            <div className="mt-3 grid grid-cols-2 gap-3 text-sm">
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Quantity
                </p>
                <p className="mt-1">{position.quantity}</p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Avg Entry
                </p>
                <p className="mt-1">{position.average_entry_price}</p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Realized PnL
                </p>
                <p className="mt-1">{position.realized_pnl}</p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Updated
                </p>
                <p className="mt-1 text-muted-foreground">
                  {formatPortfolioDate(position.updated_at)}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="hidden overflow-x-auto rounded-lg border md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Market</TableHead>
              <TableHead>Outcome</TableHead>
              <TableHead>Quantity</TableHead>
              <TableHead>Avg Entry</TableHead>
              <TableHead>Realized PnL</TableHead>
              <TableHead>Updated</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {positions.map((position) => (
              <TableRow key={position.id}>
                <TableCell className="font-mono text-xs" title={position.market_id}>
                  {truncatePortfolioId(position.market_id)}
                </TableCell>
                <TableCell>{position.outcome}</TableCell>
                <TableCell>{position.quantity}</TableCell>
                <TableCell>{position.average_entry_price}</TableCell>
                <TableCell>{position.realized_pnl}</TableCell>
                <TableCell className="text-muted-foreground">
                  {formatPortfolioDate(position.updated_at)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </>
  )
}

function SettlementsTable({ settlements }: { settlements: Settlement[] }) {
  if (settlements.length === 0) {
    return <p className="text-sm text-muted-foreground">No settlements yet.</p>
  }

  return (
    <>
      <div className="grid gap-3 md:hidden">
        {settlements.map((settlement) => (
          <div key={settlement.id} className="rounded-lg border bg-card/50 p-3">
            <div className="flex items-start justify-between gap-3">
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Market
                </p>
                <p className="mt-1 font-mono text-xs" title={settlement.market_id}>
                  {truncatePortfolioId(settlement.market_id)}
                </p>
              </div>
              <div className="text-right">
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Status
                </p>
                <p className="mt-1 text-sm font-medium">{settlement.status}</p>
              </div>
            </div>
            <div className="mt-3 grid grid-cols-2 gap-3 text-sm">
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Outcome
                </p>
                <p className="mt-1">{settlement.outcome ?? "-"}</p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Amount
                </p>
                <p className="mt-1">{settlement.amount}</p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Tx Hash
                </p>
                <p className="mt-1 font-mono text-xs">
                  {settlement.tx_hash ? truncatePortfolioId(settlement.tx_hash) : "-"}
                </p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase text-muted-foreground">
                  Settled
                </p>
                <p className="mt-1 text-muted-foreground">
                  {formatPortfolioDate(settlement.settled_at)}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="hidden overflow-x-auto rounded-lg border md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Market</TableHead>
              <TableHead>Outcome</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Tx Hash</TableHead>
              <TableHead>Settled</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {settlements.map((settlement) => (
              <TableRow key={settlement.id}>
                <TableCell className="font-mono text-xs" title={settlement.market_id}>
                  {truncatePortfolioId(settlement.market_id)}
                </TableCell>
                <TableCell>{settlement.outcome ?? "-"}</TableCell>
                <TableCell>{settlement.amount}</TableCell>
                <TableCell>{settlement.status}</TableCell>
                <TableCell className="font-mono text-xs">
                  {settlement.tx_hash ? truncatePortfolioId(settlement.tx_hash) : "-"}
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {formatPortfolioDate(settlement.settled_at)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </>
  )
}

export function PortfolioPositionCard({
  positions,
  settlements,
}: {
  positions: Position[]
  settlements: Settlement[]
}) {
  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Positions</CardTitle>
          <CardDescription>Your current market positions.</CardDescription>
        </CardHeader>
        <CardContent>
          <PositionsTable positions={positions} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Settlements</CardTitle>
          <CardDescription>
            Settlement history for resolved markets.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SettlementsTable settlements={settlements} />
        </CardContent>
      </Card>
    </div>
  )
}
