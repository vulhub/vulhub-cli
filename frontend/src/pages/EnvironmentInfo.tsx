import { useParams, Link } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Skeleton } from '@/components/ui/skeleton'
import { ScrollArea } from '@/components/ui/scroll-area'
import { ActionButtons } from '@/components/environments/ActionButtons'
import { ContainerStatus } from '@/components/environments/ContainerStatus'
import { useEnvironmentInfo, useEnvironmentStatus } from '@/hooks/useEnvironments'
import { useState } from 'react'

export function EnvironmentInfo() {
  const { '*': path } = useParams()
  const decodedPath = path ? decodeURIComponent(path) : ''

  const [tab, setTab] = useState('readme')

  const { data: info, isLoading: infoLoading } = useEnvironmentInfo(decodedPath)
  const { data: status, isLoading: statusLoading } =
    useEnvironmentStatus(decodedPath)

  if (infoLoading) {
    return (
      <div className="p-6">
        <Skeleton className="mb-4 h-8 w-48" />
        <Skeleton className="mb-8 h-4 w-96" />
        <Skeleton className="h-96" />
      </div>
    )
  }

  if (!info) {
    return (
      <div className="p-6">
        <div className="flex h-40 items-center justify-center rounded-lg border border-dashed">
          <p className="text-muted-foreground">Environment not found</p>
        </div>
      </div>
    )
  }

  const isRunning = status?.running ?? false
  const isDownloaded = info.downloaded

  return (
    <div className="p-6">
      {/* Back button */}
      <Link to="/environments">
        <Button variant="ghost" size="sm" className="mb-4">
          <ArrowLeft className="h-4 w-4" />
          Back to environments
        </Button>
      </Link>

      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold">{info.environment.name}</h1>
            <p className="text-muted-foreground">{info.environment.path}</p>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant={isRunning ? 'success' : isDownloaded ? 'secondary' : 'outline'}>
              {isRunning ? 'Running' : isDownloaded ? 'Downloaded' : 'Available'}
            </Badge>
          </div>
        </div>

        {/* Tags */}
        <div className="mt-4 flex flex-wrap gap-2">
          {info.environment.cve?.map((cve) => (
            <Badge key={cve} variant="destructive">
              {cve}
            </Badge>
          ))}
          {info.environment.tags?.map((tag) => (
            <Badge key={tag} variant="outline">
              {tag}
            </Badge>
          ))}
          {info.environment.app && (
            <Badge variant="secondary">App: {info.environment.app}</Badge>
          )}
        </div>

        {/* Actions */}
        <div className="mt-6">
          <ActionButtons
            path={decodedPath}
            running={isRunning}
            downloaded={isDownloaded}
          />
        </div>
      </div>

      {/* Container Status (if running) */}
      {isDownloaded && (
        <div className="mb-8">
          {statusLoading ? (
            <Skeleton className="h-48" />
          ) : (
            <ContainerStatus containers={status?.containers ?? []} />
          )}
        </div>
      )}

      {/* Content Tabs */}
      <Tabs value={tab} onValueChange={setTab}>
        <TabsList>
          <TabsTrigger value="readme">README</TabsTrigger>
          <TabsTrigger value="compose">docker-compose.yml</TabsTrigger>
        </TabsList>

        <TabsContent value="readme">
          <Card>
            <CardHeader>
              <CardTitle>README</CardTitle>
            </CardHeader>
            <CardContent>
              {info.readme ? (
                <ScrollArea className="h-[600px]">
                  <div className="prose prose-sm max-w-none dark:prose-invert">
                    <ReactMarkdown
                      components={{
                        code({ className, children, ...props }) {
                          const match = /language-(\w+)/.exec(className || '')
                          const isInline = !match
                          return isInline ? (
                            <code className={className} {...props}>
                              {children}
                            </code>
                          ) : (
                            <SyntaxHighlighter
                              style={vscDarkPlus}
                              language={match[1]}
                              PreTag="div"
                            >
                              {String(children).replace(/\n$/, '')}
                            </SyntaxHighlighter>
                          )
                        },
                      }}
                    >
                      {info.readme}
                    </ReactMarkdown>
                  </div>
                </ScrollArea>
              ) : (
                <p className="text-muted-foreground">No README available</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="compose">
          <Card>
            <CardHeader>
              <CardTitle>docker-compose.yml</CardTitle>
            </CardHeader>
            <CardContent>
              {info.compose_file ? (
                <ScrollArea className="h-[600px]">
                  <SyntaxHighlighter
                    style={vscDarkPlus}
                    language="yaml"
                    showLineNumbers
                  >
                    {info.compose_file}
                  </SyntaxHighlighter>
                </ScrollArea>
              ) : (
                <p className="text-muted-foreground">
                  Compose file not available. Start the environment to download
                  it.
                </p>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
