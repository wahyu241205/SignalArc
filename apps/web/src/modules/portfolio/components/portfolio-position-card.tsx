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
    <div className="overflow-x-auto rounded-lg border">
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
  )
}

function SettlementsTable({ settlements }: { settlements: Settlement[] }) {
  if (settlements.length === 0) {
    return <p className="text-sm text-muted-foreground">No settlements yet.</p>
  }

  return (
    <div className="overflow-x-auto rounded-lg border">
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
