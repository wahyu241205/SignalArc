export type AnalyticsMetric = {
  label: string
  value: string
  description: string
  unit: string
  featured: boolean
}

export type AnalyticsFactorySnapshot = {
  name: string
  address: string
  explorerUrl: string
  deploymentTx: string
  deploymentTxUrl: string
  deployer: string
  deploymentBlock: string
  deploymentTimestamp: string
  totalTransactions: number
  latestMarket: string
}

export type AnalyticsActivitySnapshot = {
  timestamp: string
  tx: string
  txUrl: string
  market: string
}

export type AnalyticsTopMarket = {
  address: string
  question: string
  collateral: string
  positionEvents: number
  explorerUrl: string
}

export type AnalyticsPublicLink = {
  label: string
  href: string
}
