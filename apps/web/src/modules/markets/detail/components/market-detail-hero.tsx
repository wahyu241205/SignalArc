import type { Market } from "@/lib/api"

export function MarketDetailHero({ market }: { market: Market }) {
  if (market.cover_image_url) {
    return (
      // Plain img is intentional for v1 user-provided remote URLs.
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={market.cover_image_url}
        alt={market.title}
        className="h-48 w-full rounded-xl object-cover sm:h-64"
      />
    )
  }

  return <div className="h-48 w-full rounded-xl bg-muted sm:h-64" aria-hidden="true" />
}
