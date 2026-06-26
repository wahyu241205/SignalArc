import { Card, CardContent, CardHeader } from "@/components/ui/card"

function SkeletonBlock({ className }: { className?: string }) {
  return <div className={`rounded bg-muted ${className ?? ""}`} />
}

export function MarketDetailLoadingSkeleton() {
  return (
    <div className="grid gap-6 lg:grid-cols-3">
      <div className="space-y-6 lg:col-span-2">
        <Card className="animate-pulse">
          <CardHeader className="space-y-4">
            <SkeletonBlock className="h-48 w-full rounded-xl sm:h-64" />
            <div className="flex items-center gap-2">
              <SkeletonBlock className="h-5 w-16" />
              <SkeletonBlock className="h-5 w-24" />
            </div>
            <SkeletonBlock className="h-8 w-3/4" />
            <SkeletonBlock className="h-4 w-full max-w-md" />
          </CardHeader>
        </Card>

        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-4 w-40" />
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-3">
              <SkeletonBlock className="h-20" />
              <SkeletonBlock className="h-20" />
            </div>
          </CardContent>
        </Card>

        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-4 w-36" />
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {Array.from({ length: 6 }).map((_, i) => (
                <SkeletonBlock key={i} className="h-10" />
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="space-y-6 lg:col-span-1">
        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-5 w-32" />
          </CardHeader>
          <CardContent className="space-y-3">
            <SkeletonBlock className="h-10" />
            <SkeletonBlock className="h-10" />
            <SkeletonBlock className="h-10" />
          </CardContent>
        </Card>

        <Card className="animate-pulse">
          <CardHeader>
            <SkeletonBlock className="h-5 w-44" />
          </CardHeader>
          <CardContent>
            <SkeletonBlock className="h-24" />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
