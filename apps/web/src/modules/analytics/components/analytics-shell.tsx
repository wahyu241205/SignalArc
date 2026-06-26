import type { ComponentProps, ReactNode } from "react"
import {
  ArrowUpRight,
  Bot,
  Check,
  CircleDollarSign,
  Database,
  ExternalLink,
  Factory,
  FileCheck2,
  Network,
  ShieldCheck,
  WalletCards,
} from "lucide-react"

import { SiteFooter } from "@/components/layout/site-footer"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

import {
  analyticsAgentIntegrationChecklist,
  analyticsBackendMetrics,
  analyticsFactory,
  analyticsFactoryProofPoints,
  analyticsLatestActivity,
  analyticsLifecycleMetrics,
  analyticsLimitations,
  analyticsMetrics,
  analyticsPublicLinks,
  analyticsStatusBadges,
  analyticsTopMarkets,
} from "../analytics-utils"
import { shortenAnalyticsAddress } from "../format"

import { AnalyticsDisclaimerCard } from "./analytics-disclaimer-card"
import { AnalyticsSummaryGrid } from "./analytics-summary-grid"

type ButtonVariant = ComponentProps<typeof Button>["variant"]

function SectionHeading({ title, description }: { title: string; description: string }) {
  return (
    <div className="flex max-w-3xl flex-col gap-2">
      <h2 className="text-xl font-semibold tracking-tight sm:text-2xl">{title}</h2>
      <p className="text-sm leading-relaxed text-muted-foreground">{description}</p>
    </div>
  )
}

function ExternalButton({
  href,
  children,
  variant = "outline",
}: {
  href: string
  children: ReactNode
  variant?: ButtonVariant
}) {
  return (
    <Button asChild variant={variant}>
      <a href={href} target="_blank" rel="noreferrer">
        {children}
        <ExternalLink data-icon="inline-end" />
      </a>
    </Button>
  )
}

function DataRow({
  label,
  value,
  mono = false,
}: {
  label: string
  value: ReactNode
  mono?: boolean
}) {
  return (
    <div className="grid gap-1 border-b border-border/40 py-3 last:border-b-0 sm:grid-cols-[180px_1fr] sm:gap-5">
      <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground">{label}</dt>
      <dd className={mono ? "min-w-0 break-all font-mono text-xs text-foreground" : "text-sm text-foreground"}>
        {value}
      </dd>
    </div>
  )
}

function AnalyticsHero() {
  return (
    <section className="grid gap-10 lg:grid-cols-[minmax(0,1.3fr)_minmax(320px,0.7fr)] lg:items-end">
      <div className="flex flex-col gap-6">
        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium text-indigo-300">Arc Testnet Proof-of-Activity Dashboard</p>
          <h1 className="max-w-4xl text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
            SignalArc{" "}
            <span className="bg-gradient-to-r from-indigo-300 via-purple-300 to-indigo-400 bg-clip-text text-transparent">
              Analytics
            </span>
          </h1>
        </div>

        <p className="max-w-3xl text-base leading-relaxed text-muted-foreground sm:text-lg">
          SignalArc publishes this dashboard as a transparent proof-of-activity layer for its Arc Testnet deployment.
          The dashboard summarizes verified factory deployment, created YES/NO market contracts, market-level position
          activity, testnet USDC collateral movement, lifecycle events, and Circle Agent Wallet integration readiness.
        </p>

        <div className="flex flex-wrap gap-2">
          {analyticsStatusBadges.map((badge, index) => (
            <Badge
              key={badge}
              variant="outline"
              className={
                index < 3
                  ? "border-indigo-500/30 bg-indigo-500/10 text-indigo-200"
                  : "border-border/70 bg-card/50 text-muted-foreground"
              }
            >
              {badge}
            </Badge>
          ))}
        </div>
      </div>

      <Card className="border-indigo-500/20 bg-card/70 shadow-2xl shadow-indigo-950/20 backdrop-blur">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-300 ring-1 ring-indigo-500/20">
              <Database className="size-4" aria-hidden="true" />
            </div>
            <div>
              <CardTitle>Data provenance</CardTitle>
              <CardDescription>Static public analytics snapshot</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <p className="text-sm leading-relaxed text-muted-foreground">
            Arc Testnet contract-level analytics are not yet reliably queryable through Dune for SignalArc&apos;s
            required metrics, so this page provides a self-hosted public analytics view backed by Arcscan-derived
            contract data and SignalArc backend references.
          </p>
          <div className="grid grid-cols-2 gap-3 border-t border-border/50 pt-4">
            <div>
              <p className="text-2xl font-semibold text-foreground">126</p>
              <p className="text-xs text-muted-foreground">Created markets</p>
            </div>
            <div>
              <p className="text-2xl font-semibold text-foreground">149.77</p>
              <p className="text-xs text-muted-foreground">Testnet USDC</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </section>
  )
}

function ExecutiveMetricsSection() {
  return (
    <section className="flex flex-col gap-7">
      <SectionHeading
        title="Executive Metrics"
        description="A contract-derived summary of verified deployment, created YES/NO markets, participating wallets, testnet collateral movement, and market lifecycle activity."
      />
      <AnalyticsSummaryGrid metrics={analyticsMetrics} />
    </section>
  )
}

function FactorySection() {
  return (
    <section className="flex flex-col gap-7">
      <SectionHeading
        title="Legacy Verified Factory Snapshot"
        description="SignalArc's historical Arc Testnet proof-of-activity is anchored by a legacy verified YES/NO market factory contract. The active factory is tracked separately; legacy metrics are kept attached to the factory that produced them."
      />

      <div className="grid gap-5 lg:grid-cols-[minmax(0,1.25fr)_minmax(280px,0.75fr)]">
        <Card className="border-indigo-500/20 bg-card/70">
          <CardHeader>
            <div className="flex items-start justify-between gap-4">
              <div className="flex items-center gap-3">
                <div className="flex size-10 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-300 ring-1 ring-indigo-500/20">
                  <Factory className="size-5" aria-hidden="true" />
                </div>
                <div>
                  <CardTitle>{analyticsFactory.name}</CardTitle>
                  <CardDescription>Canonical market deployment registry</CardDescription>
                </div>
              </div>
              <Badge variant="outline" className="border-emerald-500/30 bg-emerald-500/10 text-emerald-300">
                <ShieldCheck data-icon="inline-start" />
                Verified
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <dl>
              <DataRow label="Contract name" value={analyticsFactory.name} />
              <DataRow label="Address" value={analyticsFactory.address} mono />
              <DataRow label="Verified" value="Yes" />
              <DataRow label="Network" value="Arc Testnet" />
              <DataRow label="Deployment tx" value={analyticsFactory.deploymentTx} mono />
              <DataRow label="Deployer" value={analyticsFactory.deployer} mono />
              <DataRow label="Deployment block" value={analyticsFactory.deploymentBlock} />
              <DataRow label="Deployment timestamp" value={analyticsFactory.deploymentTimestamp} mono />
              <DataRow label="Factory transactions" value={analyticsFactory.totalTransactions} />
              <DataRow label="Latest market created" value={analyticsFactory.latestMarket} mono />
            </dl>
          </CardContent>
        </Card>

        <div className="flex flex-col gap-5">
          <Card className="border-border/60 bg-card/60">
            <CardHeader>
              <CardTitle>Factory proof points</CardTitle>
              <CardDescription>Explorer-verifiable deployment references</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-4">
              {analyticsFactoryProofPoints.map((item) => (
                <div key={item.label} className="flex items-center gap-3">
                  <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-muted text-indigo-300">
                    <Factory className="size-4" aria-hidden="true" />
                  </div>
                  <div>
                    <p className="text-sm font-medium">{item.value}</p>
                    <p className="text-xs text-muted-foreground">{item.label}</p>
                  </div>
                </div>
              ))}
            </CardContent>
          </Card>

          <div className="flex flex-col gap-2">
            <ExternalButton href={analyticsFactory.explorerUrl}>View Legacy Factory on Arcscan</ExternalButton>
            <ExternalButton href={analyticsFactory.deploymentTxUrl}>View Legacy Deployment Transaction</ExternalButton>
            <ExternalButton href={analyticsLatestActivity.txUrl}>View Legacy Latest Activity</ExternalButton>
          </div>
        </div>
      </div>
    </section>
  )
}

function MarketsSection() {
  return (
    <section className="flex flex-col gap-7">
      <SectionHeading
        title="Created YES/NO Markets"
        description="Each market created by the factory represents an individual YES/NO prediction market contract. Trading and lifecycle activity lives at the market-contract level, while the factory provides the canonical market creation registry."
      />

      <Card className="border-border/60 bg-card/60">
        <CardHeader className="border-b border-border/50">
          <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <CardTitle>Top markets by collateral</CardTitle>
              <CardDescription>Five leading contracts from the 126-market Arc Testnet deployment</CardDescription>
            </div>
            <Badge variant="outline">Testnet USDC</Badge>
          </div>
        </CardHeader>
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="min-w-[320px] pl-4">Question</TableHead>
                <TableHead>Contract</TableHead>
                <TableHead className="text-right">Collateral</TableHead>
                <TableHead className="text-right">Position events</TableHead>
                <TableHead className="pr-4 text-right">Explorer</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {analyticsTopMarkets.map((market) => (
                <TableRow key={market.address}>
                  <TableCell className="max-w-[520px] whitespace-normal pl-4 font-medium">{market.question}</TableCell>
                  <TableCell className="font-mono text-xs text-muted-foreground" title={market.address}>
                    {shortenAnalyticsAddress(market.address)}
                  </TableCell>
                  <TableCell className="text-right font-medium">
                    {market.collateral} <span className="text-xs text-muted-foreground">testnet USDC</span>
                  </TableCell>
                  <TableCell className="text-right">{market.positionEvents}</TableCell>
                  <TableCell className="pr-4 text-right">
                    <Button asChild variant="ghost" size="sm">
                      <a
                        href={market.explorerUrl}
                        target="_blank"
                        rel="noreferrer"
                        aria-label={`View ${market.address} on Arcscan`}
                      >
                        Arcscan
                        <ArrowUpRight data-icon="inline-end" />
                      </a>
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </section>
  )
}

function LifecycleSection() {
  return (
    <section className="flex flex-col gap-7">
      <SectionHeading
        title="Lifecycle Activity"
        description="SignalArc's Arc Testnet deployment includes market creation, position activity, market resolution, cancellation, and claim/refund lifecycle events. These metrics demonstrate end-to-end contract lifecycle coverage across the testnet deployment."
      />

      <div className="grid gap-5 lg:grid-cols-[0.8fr_1.2fr]">
        <div className="grid grid-cols-3 gap-3 self-start">
          {analyticsLifecycleMetrics.map((item) => (
            <Card key={item.label} size="sm" className="border-border/60 bg-card/60">
              <CardHeader>
                <CardTitle className="text-3xl font-semibold">{item.value}</CardTitle>
                <CardDescription className="text-xs">{item.label}</CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>

        <Card className="border-indigo-500/20 bg-gradient-to-r from-card to-indigo-500/5">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex size-9 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-300">
                <Network className="size-4" aria-hidden="true" />
              </div>
              <div>
                <CardTitle>Latest recorded activity</CardTitle>
                <CardDescription>Most recent market creation in the analytics snapshot</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <dl>
              <DataRow label="Timestamp" value={analyticsLatestActivity.timestamp} mono />
              <DataRow label="Transaction" value={analyticsLatestActivity.tx} mono />
              <DataRow label="Market" value={analyticsLatestActivity.market} mono />
            </dl>
            <div className="pt-4">
              <ExternalButton href={analyticsLatestActivity.txUrl} variant="default">
                Inspect Latest Activity
              </ExternalButton>
            </div>
          </CardContent>
        </Card>
      </div>
    </section>
  )
}

function AgentWalletSection() {
  return (
    <section className="flex flex-col gap-7">
      <SectionHeading
        title="Circle Agent Wallet Integration"
        description="SignalArc integrates Circle Agent Wallet / Programmable Wallet flows in the backend for agent wallet onboarding, OTP verification, wallet resolution, session activation, wallet balance lookup, faucet requests, and backend-driven smart contract execution on Arc Testnet."
      />

      <div className="grid gap-5 lg:grid-cols-2">
        <Card className="border-indigo-500/20 bg-card/70">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex size-10 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-300">
                <Bot className="size-5" aria-hidden="true" />
              </div>
              <div>
                <CardTitle>Integration readiness</CardTitle>
                <CardDescription>Backend and agent execution surfaces</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="grid gap-3 sm:grid-cols-2">
            {analyticsAgentIntegrationChecklist.map((item) => (
              <div key={item} className="flex items-start gap-2 rounded-lg border border-border/50 bg-background/30 p-3">
                <div className="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-indigo-500/10 text-indigo-300">
                  <Check className="size-3" aria-hidden="true" />
                </div>
                <span className="text-sm leading-snug">{item}</span>
              </div>
            ))}
          </CardContent>
        </Card>

        <div className="flex flex-col gap-5">
          <Card className="border-amber-500/20 bg-amber-500/5">
            <CardHeader>
              <div className="flex items-center gap-3">
                <WalletCards className="size-5 text-amber-300" aria-hidden="true" />
                <CardTitle>Attribution boundary</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                Direct Circle Agent Wallet attribution is not visible from chain data alone. On-chain analytics show
                wallet and contract interactions, while Circle Agent Wallet session, onboarding, and intent metrics are
                sourced from SignalArc backend data.
              </p>
            </CardContent>
          </Card>

          <Card className="border-border/60 bg-card/60">
            <CardHeader>
              <CardTitle>Backend-sourced metrics to add</CardTitle>
              <CardDescription>
                These operational metrics require backend references and are not represented as onchain totals.
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-x-5 gap-y-2 sm:grid-cols-2">
              {analyticsBackendMetrics.map((metric) => (
                <div key={metric} className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="size-1.5 rounded-full bg-indigo-400/70" />
                  {metric}
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  )
}

function PublicLinksSection() {
  return (
    <section className="grid gap-5 lg:grid-cols-[1.1fr_0.9fr]">
      <Card className="border-border/60 bg-card/60">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-lg bg-muted text-indigo-300">
              <ArrowUpRight className="size-4" aria-hidden="true" />
            </div>
            <div>
              <CardTitle>Public Links</CardTitle>
              <CardDescription>Product, documentation, source, and explorer references</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="grid gap-2 sm:grid-cols-2">
          {analyticsPublicLinks.map((link) => (
            <a
              key={link.label}
              href={link.href}
              target="_blank"
              rel="noreferrer"
              className="group flex items-center justify-between gap-3 rounded-lg border border-border/50 bg-background/30 px-3 py-3 text-sm transition-colors hover:border-indigo-500/30 hover:bg-indigo-500/5"
            >
              <span>{link.label}</span>
              <ArrowUpRight
                className="size-4 text-muted-foreground transition-colors group-hover:text-indigo-300"
                aria-hidden="true"
              />
            </a>
          ))}
        </CardContent>
      </Card>

      <AnalyticsDisclaimerCard
        title="Proof, not production claims"
        description="Arc Testnet context is preserved throughout"
        className="border-indigo-500/20 bg-gradient-to-br from-indigo-500/10 via-card to-purple-500/5"
        icon={
          <div className="flex size-9 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-300">
            <CircleDollarSign className="size-4" aria-hidden="true" />
          </div>
        }
      >
        <p className="text-sm leading-relaxed text-muted-foreground">
          The strongest signals in this snapshot are a verified factory, 126 created markets, 806 total trades, 218
          unique participating wallets, 149.77 testnet USDC in aggregate collateral movement, 105 resolved markets, and
          393 claim events.
        </p>
      </AnalyticsDisclaimerCard>
    </section>
  )
}

function NotesSection() {
  return (
    <section className="grid gap-5 lg:grid-cols-2">
      <AnalyticsDisclaimerCard
        title="Dune Availability Note"
        description="Current Arc Testnet indexing constraint"
        icon={
          <div className="flex size-9 items-center justify-center rounded-lg bg-muted text-indigo-300">
            <Database className="size-4" aria-hidden="true" />
          </div>
        }
      >
        <div className="flex flex-col gap-4 text-sm leading-relaxed text-muted-foreground">
          <p>
            SignalArc currently operates on Arc Testnet. At this stage, Arc Testnet data is not yet reliably queryable
            through Dune for the contract-level analytics required by this dashboard. For that reason, SignalArc
            publishes this public proof-of-activity page using Arcscan-derived contract data, verified factory
            references, market-level activity, backend integration references, and explorer links.
          </p>
          <p>
            A Dune dashboard or another third-party analytics dashboard will be added once Arc Testnet indexing becomes
            reliably available.
          </p>
        </div>
      </AnalyticsDisclaimerCard>

      <AnalyticsDisclaimerCard
        title="Methodology and Limitations"
        description="Transparent interpretation boundaries"
        icon={
          <div className="flex size-9 items-center justify-center rounded-lg bg-muted text-indigo-300">
            <FileCheck2 className="size-4" aria-hidden="true" />
          </div>
        }
      >
        <ul className="flex flex-col gap-3">
          {analyticsLimitations.map((limitation) => (
            <li key={limitation} className="flex items-start gap-3 text-sm leading-relaxed text-muted-foreground">
              <span className="mt-2 size-1.5 shrink-0 rounded-full bg-indigo-400/70" />
              {limitation}
            </li>
          ))}
        </ul>
      </AnalyticsDisclaimerCard>
    </section>
  )
}

export function AnalyticsShell() {
  return (
    <>
      <div className="relative overflow-hidden">
        <div className="pointer-events-none absolute inset-x-0 top-0 h-[560px] bg-[radial-gradient(circle_at_top_left,rgba(99,102,241,0.12),transparent_38%),radial-gradient(circle_at_82%_10%,rgba(168,85,247,0.1),transparent_34%)]" />

        <div className="relative px-4 py-10 sm:px-6 sm:py-14 lg:px-8">
          <div className="mx-auto flex w-full max-w-7xl flex-col gap-16 sm:gap-20">
            <AnalyticsHero />
            <ExecutiveMetricsSection />
            <Separator className="opacity-40" />
            <FactorySection />
            <MarketsSection />
            <LifecycleSection />
            <Separator className="opacity-40" />
            <AgentWalletSection />
            <PublicLinksSection />
            <NotesSection />
          </div>
        </div>
      </div>
      <SiteFooter />
    </>
  )
}
