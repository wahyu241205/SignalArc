import { Card, CardContent, CardHeader } from "@/components/ui/card"
import type { Market } from "@/lib/api"

import { arcscanContractUrl, formatMarketDate } from "../format"

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
        {label}
      </dt>
      <dd className="mt-1 text-sm font-medium text-foreground">{value}</dd>
    </div>
  )
}

export function MarketMetadataCard({ market }: { market: Market }) {
  return (
    <>
      <Card>
        <CardHeader>
          <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Market Parameters
          </h3>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <DetailItem label="Closes" value={formatMarketDate(market.closes_at)} />
            {market.opens_at ? (
              <DetailItem label="Opens" value={formatMarketDate(market.opens_at)} />
            ) : null}
            {market.resolved_at ? (
              <DetailItem
                label="Resolved"
                value={formatMarketDate(market.resolved_at)}
              />
            ) : null}
            {market.winning_outcome ? (
              <DetailItem label="Winning Outcome" value={market.winning_outcome} />
            ) : null}
            <DetailItem label="Created" value={formatMarketDate(market.created_at)} />
            <DetailItem
              label="Onchain Status"
              value={market.onchain_deployment_status}
            />
            {market.market_contract_address ? (
              <div>
                <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">
                  Market Contract
                </dt>
                <dd className="mt-1 text-sm font-medium">
                  <a
                    href={arcscanContractUrl(market.market_contract_address)}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="break-all font-mono text-indigo-400 transition-colors hover:text-indigo-300"
                  >
                    {market.market_contract_address}
                  </a>
                </dd>
              </div>
            ) : null}
          </div>
        </CardContent>
      </Card>

      {market.resolution_source ? (
        <Card>
          <CardHeader>
            <h3 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Resolution
            </h3>
          </CardHeader>
          <CardContent>
            <DetailItem
              label="Resolution Source"
              value={market.resolution_source}
            />
          </CardContent>
        </Card>
      ) : null}
    </>
  )
}
