import { Card, CardContent } from "@/components/ui/card"

export function PortfolioLoadingSkeleton() {
  return (
    <Card className="animate-pulse">
      <CardContent className="py-8">
        <div className="h-4 w-1/2 rounded bg-muted" />
      </CardContent>
    </Card>
  )
}
