"use client"

import Link from "next/link"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"

export function MarketDetailErrorState({
  message,
  requestId,
}: {
  message: string
  requestId: string | null
}) {
  return (
    <div className="mx-auto max-w-lg py-12">
      <Card className="border-destructive/30 bg-destructive/5">
        <CardHeader className="space-y-2 text-center">
          <CardTitle className="text-lg text-destructive">
            Unable to load market
          </CardTitle>
          <CardDescription className="text-sm text-muted-foreground">
            {message}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 text-center">
          {requestId ? (
            <p className="font-mono text-xs text-muted-foreground">
              Request ID: {requestId}
            </p>
          ) : null}
          <Separator />
          <div className="flex items-center justify-center gap-3">
            <Button asChild size="sm" variant="outline">
              <Link href="/markets">Back to markets</Link>
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => window.location.reload()}
            >
              Retry
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
