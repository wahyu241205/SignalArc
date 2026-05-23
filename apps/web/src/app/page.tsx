"use client"

import { ConnectButton } from "@rainbow-me/rainbowkit"
import Link from "next/link"
import { useAccount } from "wagmi"

import { SiteFooter } from "@/components/layout/site-footer"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

const features = [
  {
    href: "/markets",
    label: "Explore Markets",
    description: "Discover and trade USDC-settled prediction markets.",
    icon: (
      <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
      </svg>
    ),
  },
  {
    href: "/markets/new",
    label: "Create a Market",
    description: "Launch your own event market with custom outcomes.",
    icon: (
      <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
      </svg>
    ),
  },
  {
    href: "/portfolio",
    label: "Your Portfolio",
    description: "Track positions, outcomes, and settlement history.",
    icon: (
      <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 9m18 0V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v3" />
      </svg>
    ),
  },
  {
    href: "/intelligence",
    label: "Market Intelligence",
    description: "Real-time probability signals and market data.",
    icon: (
      <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z" />
      </svg>
    ),
  },
]

export default function Home() {
  const { isConnected } = useAccount()

  return (
    <>
      <div className="px-4 py-8 sm:px-6 lg:px-8">
        <div className="mx-auto flex w-full max-w-7xl flex-col gap-16">        {/* Hero */}
        <section className="flex flex-col items-center gap-6 pt-12 text-center lg:pt-20">
          <h1 className="max-w-4xl text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
            Trade the future with
            <span className="block bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
              prediction markets on Arc
            </span>
          </h1>
          <p className="max-w-2xl text-lg leading-relaxed text-muted-foreground">
            Create, trade, and resolve USDC-settled event markets with transparent
            settlement and real-time probability signals.
          </p>

          <div className="flex flex-col items-center gap-4 pt-4 sm:flex-row">
            {!isConnected ? (
              <ConnectButton label="Connect Wallet" />
            ) : (
              <Button asChild size="lg">
                <Link href="/markets">Explore Markets</Link>
              </Button>
            )}
            <Button asChild variant="outline" size="lg">
              <Link href="/markets/new">Create a Market</Link>
            </Button>
          </div>
        </section>

        {/* Feature cards */}
        <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {features.map((feature) => (
            <Link key={feature.href} href={feature.href} className="group">
              <Card className="h-full transition-colors group-hover:border-indigo-500/40 group-hover:bg-card/80">
                <CardHeader>
                  <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-400 transition-colors group-hover:bg-indigo-500/20">
                    {feature.icon}
                  </div>
                  <CardTitle className="text-base">{feature.label}</CardTitle>
                  <CardDescription>{feature.description}</CardDescription>
                </CardHeader>
              </Card>
            </Link>
          ))}
        </section>

        {/* Value props */}
        <section className="grid gap-4 sm:grid-cols-3">
          <Card className="border-indigo-500/10">
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-indigo-400">USDC</div>
              <p className="mt-1 text-sm text-muted-foreground">Settled in USDC on Arc</p>
            </CardContent>
          </Card>
          <Card className="border-indigo-500/10">
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-indigo-400">Binary Markets</div>
              <p className="mt-1 text-sm text-muted-foreground">YES / NO outcome resolution</p>
            </CardContent>
          </Card>
          <Card className="border-indigo-500/10">
            <CardContent className="pt-6">
              <div className="text-2xl font-bold text-indigo-400">API-First</div>
              <p className="mt-1 text-sm text-muted-foreground">Built for agents and developers</p>
            </CardContent>
          </Card>
        </section>

        {/* Network info — subtle, not hero content */}
        <section>
          <details className="group">
            <summary className="cursor-pointer text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
              Network &amp; Developer Info
            </summary>
            <div className="mt-4">
              <NetworkInfoPanel />
            </div>
          </details>
        </section>
        </div>
      </div>
      <SiteFooter />
    </>
  )
}

function NetworkInfoPanel() {
  return (
    <Card className="border-border/50">
      <CardHeader>
        <div className="flex items-center gap-2">
          <CardTitle className="text-sm">Arc Testnet</CardTitle>
          <Badge variant="outline" className="border-indigo-500/30 bg-indigo-500/10 text-indigo-300 text-xs">
            Testnet
          </Badge>
        </div>
        <CardDescription className="text-xs">
          SignalArc is deployed on Arc Testnet (Chain ID 5042002). Contract details are available via the backend API.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <dl className="grid gap-3 text-sm sm:grid-cols-3">
          <div>
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Network</dt>
            <dd className="mt-0.5 text-foreground">Arc Testnet</dd>
          </div>
          <div>
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Chain ID</dt>
            <dd className="mt-0.5 text-foreground">5042002</dd>
          </div>
          <div>
            <dt className="text-xs font-medium uppercase tracking-wider text-muted-foreground/70">Explorer</dt>
            <dd className="mt-0.5">
              <a
                href="https://testnet.arcscan.app"
                target="_blank"
                rel="noopener noreferrer"
                className="text-indigo-400 hover:text-indigo-300 transition-colors"
              >
                testnet.arcscan.app
              </a>
            </dd>
          </div>
        </dl>
      </CardContent>
    </Card>
  )
}
