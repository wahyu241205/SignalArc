"use client"

import { Label } from "@/components/ui/label"

import type { TradeOutcome } from "../types"

export function TradeSideSelector({
  value,
  onChange,
}: {
  value: TradeOutcome
  onChange: (value: TradeOutcome) => void
}) {
  return (
    <div className="grid gap-2">
      <Label htmlFor="outcome">Outcome</Label>
      <div className="grid grid-cols-2 gap-2">
        <label className="flex cursor-pointer items-center justify-center rounded-lg border border-border bg-input px-4 py-2.5 text-sm font-medium transition-colors has-[:checked]:border-green-500/50 has-[:checked]:bg-green-500/10 has-[:checked]:text-green-300">
          <input
            type="radio"
            name="outcome"
            value="YES"
            checked={value === "YES"}
            className="sr-only"
            onChange={() => onChange("YES")}
          />
          YES
        </label>
        <label className="flex cursor-pointer items-center justify-center rounded-lg border border-border bg-input px-4 py-2.5 text-sm font-medium transition-colors has-[:checked]:border-red-500/50 has-[:checked]:bg-red-500/10 has-[:checked]:text-red-300">
          <input
            type="radio"
            name="outcome"
            value="NO"
            checked={value === "NO"}
            className="sr-only"
            onChange={() => onChange("NO")}
          />
          NO
        </label>
      </div>
    </div>
  )
}
