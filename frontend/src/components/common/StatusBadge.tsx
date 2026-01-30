import { Badge } from '@/components/ui/badge'

interface StatusBadgeProps {
  running: boolean
  downloaded?: boolean
}

export function StatusBadge({ running, downloaded }: StatusBadgeProps) {
  if (running) {
    return (
      <Badge className="border-0 bg-emerald-400/20 font-mono text-[10px] font-normal text-emerald-400">
        <span className="mr-1.5 inline-block h-1.5 w-1.5 animate-pulse rounded-full bg-emerald-400" />
        LIVE
      </Badge>
    )
  }
  if (downloaded) {
    return (
      <Badge className="border-0 bg-blue-400/20 font-mono text-[10px] font-normal text-blue-400">
        READY
      </Badge>
    )
  }
  return (
    <Badge variant="outline" className="border-border/50 font-mono text-[10px] font-normal text-muted-foreground">
      AVAILABLE
    </Badge>
  )
}
