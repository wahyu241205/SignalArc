import Link from "next/link"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { ContractStatusCard } from "@/features/arc/contract-status-card"
import { localApiBaseUrl } from "@/lib/api"

const primaryLinks = [
  { href: "/markets", label: "Markets", description: "Browse local market records." },
  { href: "/markets/new", label: "Create market", description: "Create an OPEN local test market." },
  { href: "/portfolio", label: "Portfolio", description: "Load read-only positions and settlements." },
  { href: "/intelligence", label: "Intelligence", description: "Inspect agent-readable market data." },
]

export default function Home() {
  return (
    <main className="min-h-screen bg-background px-4 py-8 text-foreground sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-8">
        <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr] lg:items-start">
          <div className="space-y-5">
            <div className="space-y-3">
              <p className="text-sm font-medium text-muted-foreground">
                Local Arc Testnet prototype
              </p>
              <h1 className="max-w-3xl text-4xl font-semibold tracking-tight sm:text-5xl">
                SignalArc local MVP
              </h1>
              <p className="max-w-2xl text-base leading-7 text-muted-foreground">
                Create and inspect USDC-settled prediction market records through the
                local Go API. This browser flow is prototype-only and does not execute
                wallet transactions, Circle actions, or settlement writes.
              </p>
            </div>

            <div className="flex flex-col gap-3 sm:flex-row">
              <Button asChild>
                <Link href="/markets">Open markets</Link>
              </Button>
              <Button asChild variant="outline">
                <Link href="/markets/new">Create local market</Link>
              </Button>
            </div>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Local services</CardTitle>
              <CardDescription>
                Expected local development endpoints for this MVP.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <dl className="grid gap-4 text-sm">
                <div>
                  <dt className="font-medium text-muted-foreground">Frontend</dt>
                  <dd className="font-mono text-xs">http://localhost:3000</dd>
                </div>
                <div>
                  <dt className="font-medium text-muted-foreground">Backend API</dt>
                  <dd className="font-mono text-xs">{localApiBaseUrl}</dd>
                </div>
                <div>
                  <dt className="font-medium text-muted-foreground">Scope</dt>
                  <dd>Local browser MVP, Arc Testnet reference, live use not approved.</dd>
                </div>
              </dl>
            </CardContent>
          </Card>
        </section>

        <section className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {primaryLinks.map((item) => (
            <Card key={item.href}>
              <CardHeader>
                <CardTitle className="text-base">{item.label}</CardTitle>
                <CardDescription>{item.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <Button asChild size="sm" variant="outline">
                  <Link href={item.href}>Open</Link>
                </Button>
              </CardContent>
            </Card>
          ))}
        </section>

        <ContractStatusCard />
      </div>
    </main>
  )
}
