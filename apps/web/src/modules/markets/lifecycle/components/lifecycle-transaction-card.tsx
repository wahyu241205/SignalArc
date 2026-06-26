import type { Hash } from "viem"

import { TransactionLink } from "@/modules/wallet"

export function LifecycleTransactionLink({ hash }: { hash: Hash }) {
  return <TransactionLink hash={hash} />
}

export function LifecycleTransactionCard({ hash }: { hash: Hash }) {
  return (
    <>
      Transaction: <LifecycleTransactionLink hash={hash} />
    </>
  )
}
