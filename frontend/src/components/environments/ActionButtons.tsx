import { useState } from 'react'
import { Play, Square, RotateCcw, Trash2, Loader2, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  useStartEnvironment,
  useStopEnvironment,
  useRestartEnvironment,
  useCleanEnvironment,
} from '@/hooks/useEnvironments'

interface ActionButtonsProps {
  path: string
  running: boolean
  downloaded: boolean
}

export function ActionButtons({ path, running, downloaded }: ActionButtonsProps) {
  const [cleanDialogOpen, setCleanDialogOpen] = useState(false)

  const startMutation = useStartEnvironment()
  const stopMutation = useStopEnvironment()
  const restartMutation = useRestartEnvironment()
  const cleanMutation = useCleanEnvironment()

  const isLoading =
    startMutation.isPending ||
    stopMutation.isPending ||
    restartMutation.isPending ||
    cleanMutation.isPending

  const handleStart = () => {
    startMutation.mutate({ path })
  }

  const handleStop = () => {
    stopMutation.mutate(path)
  }

  const handleRestart = () => {
    restartMutation.mutate(path)
  }

  const handleClean = () => {
    cleanMutation.mutate(
      { path, options: { remove_volumes: true, remove_files: true } },
      { onSuccess: () => setCleanDialogOpen(false) }
    )
  }

  return (
    <>
      <div className="flex flex-wrap gap-2">
        {running ? (
          <>
            <Button
              variant="outline"
              onClick={handleStop}
              disabled={isLoading}
              className="border-amber-400/50 font-mono text-xs uppercase text-amber-400 hover:bg-amber-400/20 hover:text-amber-400"
            >
              {stopMutation.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Square className="h-4 w-4" />
              )}
              STOP
            </Button>
            <Button
              variant="outline"
              onClick={handleRestart}
              disabled={isLoading}
              className="border-cyan-400/50 font-mono text-xs uppercase text-cyan-400 hover:bg-cyan-400/20 hover:text-cyan-400"
            >
              {restartMutation.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <RotateCcw className="h-4 w-4" />
              )}
              RESTART
            </Button>
          </>
        ) : (
          <Button
            onClick={handleStart}
            disabled={isLoading}
            className="bg-emerald-500 font-mono text-xs uppercase text-white hover:bg-emerald-400"
          >
            {startMutation.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Play className="h-4 w-4" />
            )}
            {downloaded ? 'START' : 'DEPLOY'}
          </Button>
        )}
        {downloaded && (
          <Button
            variant="outline"
            onClick={() => setCleanDialogOpen(true)}
            disabled={isLoading}
            className="border-red-400/50 font-mono text-xs uppercase text-red-400 hover:bg-red-400/20 hover:text-red-400"
          >
            <Trash2 className="h-4 w-4" />
            CLEAN
          </Button>
        )}
      </div>

      <Dialog open={cleanDialogOpen} onOpenChange={setCleanDialogOpen}>
        <DialogContent className="border-border/50 bg-card">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 font-mono text-lg uppercase tracking-wider text-red-400">
              <AlertTriangle className="h-5 w-5" />
              Clean Environment
            </DialogTitle>
            <DialogDescription className="font-mono text-sm text-muted-foreground">
              This will permanently remove containers, volumes, and local files.
              This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setCleanDialogOpen(false)}
              disabled={cleanMutation.isPending}
              className="font-mono text-xs uppercase"
            >
              CANCEL
            </Button>
            <Button
              onClick={handleClean}
              disabled={cleanMutation.isPending}
              className="bg-red-500 font-mono text-xs uppercase text-white hover:bg-red-600"
            >
              {cleanMutation.isPending && (
                <Loader2 className="h-4 w-4 animate-spin" />
              )}
              CONFIRM CLEAN
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
