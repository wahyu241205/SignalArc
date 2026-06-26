"use client"

import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export function TradeAmountInput({
  value,
  onChange,
}: {
  value: string
  onChange: (value: string) => void
}) {
  return (
    <div className="grid gap-2">
      <Label htmlFor="amount">Amount (USDC)</Label>
      <Input
        id="amount"
        inputMode="decimal"
        min="0"
        name="amount"
        onChange={(event) => onChange(event.target.value)}
        required
        step="0.000001"
        value={value}
      />
    </div>
  )
}
