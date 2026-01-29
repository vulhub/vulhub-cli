import { useState } from 'react'
import { Play, Square, RotateCcw, Trash2, Loader2 } from 'lucide-react'
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
            >
              {stopMutation.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Square className="h-4 w-4" />
              )}
              Stop
            </Button>
            <Button
              variant="outline"
              onClick={handleRestart}
              disabled={isLoading}
            >
              {restartMutation.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <RotateCcw className="h-4 w-4" />
              )}
              Restart
            </Button>
          </>
        ) : (
          <Button onClick={handleStart} disabled={isLoading}>
            {startMutation.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Play className="h-4 w-4" />
            )}
            Start
          </Button>
        )}
        {downloaded && (
          <Button
            variant="destructive"
            onClick={() => setCleanDialogOpen(true)}
            disabled={isLoading}
          >
            <Trash2 className="h-4 w-4" />
            Clean
          </Button>
        )}
      </div>

      <Dialog open={cleanDialogOpen} onOpenChange={setCleanDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Clean Environment</DialogTitle>
            <DialogDescription>
              Are you sure you want to clean this environment? This will remove
              containers, volumes, and local files.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setCleanDialogOpen(false)}
              disabled={cleanMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleClean}
              disabled={cleanMutation.isPending}
            >
              {cleanMutation.isPending && (
                <Loader2 className="h-4 w-4 animate-spin" />
              )}
              Clean
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
