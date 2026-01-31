import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, FileText, FileCode, AlertTriangle } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeRaw from 'rehype-raw'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Skeleton } from '@/components/ui/skeleton'
import { ActionButtons } from '@/components/environments/ActionButtons'
import { ContainerStatus } from '@/components/environments/ContainerStatus'
import { useEnvironmentInfo, useEnvironmentStatus } from '@/hooks/useEnvironments'
import React, { useState } from 'react'

// GitHub repository base URLs for vulhub
const GITHUB_RAW_BASE = 'https://raw.githubusercontent.com/vulhub/vulhub/master'
const GITHUB_BLOB_BASE = 'https://github.com/vulhub/vulhub/blob/master'

/**
 * Check if a URL is a relative path (not absolute or protocol-relative)
 */
function isRelativePath(url: string | undefined): boolean {
  if (!url) return false
  // Absolute URLs or protocol-relative URLs
  if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('//')) {
    return false
  }
  // Data URLs, mailto, tel, etc.
  if (url.includes(':')) {
    return false
  }
  // Anchor links
  if (url.startsWith('#')) {
    return false
  }
  return true
}

/**
 * Resolve a relative path to an absolute GitHub URL
 * @param relativePath - The relative path from the markdown
 * @param environmentPath - The environment path (e.g., "log4j/CVE-2021-44228")
 * @param baseUrl - The GitHub base URL to use
 */
function resolveGitHubUrl(
  relativePath: string,
  environmentPath: string,
  baseUrl: string
): string {
  // Remove leading "./" if present
  let cleanPath = relativePath.replace(/^\.\//, '')
  
  // Handle parent directory references "../"
  const envParts = environmentPath.split('/')
  while (cleanPath.startsWith('../')) {
    cleanPath = cleanPath.slice(3)
    envParts.pop()
  }
  
  const resolvedEnvPath = envParts.join('/')
  
  if (resolvedEnvPath) {
    return `${baseUrl}/${resolvedEnvPath}/${cleanPath}`
  }
  return `${baseUrl}/${cleanPath}`
}

export function EnvironmentInfo() {
  const { '*': path } = useParams()
  const decodedPath = path ? decodeURIComponent(path) : ''

  const [tab, setTab] = useState('readme')

  const { data: info, isLoading: infoLoading } = useEnvironmentInfo(decodedPath)
  
  // Only poll status when environment is downloaded
  const isDownloaded = info?.downloaded ?? false
  const { data: status, isLoading: statusLoading } =
    useEnvironmentStatus(decodedPath, { enabled: isDownloaded })

  if (infoLoading) {
    return (
      <div className="relative min-h-full p-6">
        <div className="cyber-grid pointer-events-none absolute inset-0 opacity-50" />
        <div className="relative">
          <Skeleton className="mb-4 h-8 w-48" />
          <Skeleton className="mb-8 h-4 w-96" />
          <Skeleton className="h-96" />
        </div>
      </div>
    )
  }

  if (!info) {
    return (
      <div className="relative min-h-full p-6">
        <div className="cyber-grid pointer-events-none absolute inset-0 opacity-50" />
        <div className="relative">
          <div className="flex h-40 flex-col items-center justify-center rounded-lg border border-dashed border-border/50 bg-card/30">
            <AlertTriangle className="mb-2 h-8 w-8 text-amber-400" />
            <p className="font-mono text-sm text-muted-foreground">
              Environment not found
            </p>
          </div>
        </div>
      </div>
    )
  }

  const isRunning = status?.running ?? false

  return (
    <div className="relative min-h-full p-6">
      {/* Background grid */}
      <div className="cyber-grid pointer-events-none absolute inset-0 opacity-50" />

      <div className="relative">
        {/* Back button */}
        <Link to="/environments">
          <Button
            variant="ghost"
            size="sm"
            className="mb-4 font-mono text-xs text-muted-foreground hover:text-primary"
          >
            <ArrowLeft className="h-4 w-4" />
            BACK
          </Button>
        </Link>

        {/* Header */}
        <div className="mb-8">
          <div className="flex items-start justify-between">
            <div>
              <h1 className="font-mono text-2xl font-bold tracking-tight text-foreground">
                {info.environment.name}
              </h1>
              <p className="font-mono text-sm text-muted-foreground">
                {info.environment.path}
              </p>
            </div>
            <div className="flex items-center gap-2">
              {isRunning ? (
                <Badge className="border-0 bg-emerald-400/20 font-mono text-xs font-normal text-emerald-400">
                  <span className="mr-1.5 inline-block h-1.5 w-1.5 animate-pulse rounded-full bg-emerald-400" />
                  LIVE
                </Badge>
              ) : isDownloaded ? (
                <Badge className="border-0 bg-blue-400/20 font-mono text-xs font-normal text-blue-400">
                  READY
                </Badge>
              ) : (
                <Badge
                  variant="outline"
                  className="border-border/50 font-mono text-xs font-normal text-muted-foreground"
                >
                  AVAILABLE
                </Badge>
              )}
            </div>
          </div>

          {/* Tags */}
          <div className="mt-4 flex flex-wrap gap-2">
            {info.environment.cve?.map((cve) => (
              <Badge
                key={cve}
                className="border-0 bg-red-500/20 font-mono text-xs font-normal text-red-400"
              >
                {cve}
              </Badge>
            ))}
            {info.environment.tags?.map((tag) => (
              <Badge
                key={tag}
                variant="outline"
                className="border-border/50 font-mono text-xs font-normal text-muted-foreground"
              >
                {tag}
              </Badge>
            ))}
            {info.environment.app && (
              <Badge className="border-0 bg-cyan-400/20 font-mono text-xs font-normal text-cyan-400">
                APP: {info.environment.app}
              </Badge>
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
          <TabsList className="border border-border/50 bg-card/80">
            <TabsTrigger
              value="readme"
              className="flex items-center gap-2 font-mono text-xs uppercase data-[state=active]:bg-primary/20 data-[state=active]:text-primary"
            >
              <FileText className="h-4 w-4" />
              README
            </TabsTrigger>
            <TabsTrigger
              value="compose"
              className="flex items-center gap-2 font-mono text-xs uppercase data-[state=active]:bg-primary/20 data-[state=active]:text-primary"
            >
              <FileCode className="h-4 w-4" />
              COMPOSE
            </TabsTrigger>
          </TabsList>

          <TabsContent value="readme">
            <Card className="border-border/50 bg-card/80 backdrop-blur">
              <CardHeader className="border-b border-border/50">
                <CardTitle className="flex items-center gap-2 font-mono text-sm uppercase tracking-wider">
                  <FileText className="h-4 w-4 text-primary" />
                  README.md
                </CardTitle>
              </CardHeader>
              <CardContent className="pt-6">
                {info.readme ? (
                  <article className="prose prose-lg max-w-none">
                      <ReactMarkdown
                        remarkPlugins={[remarkGfm]}
                        rehypePlugins={[rehypeRaw]}
                        components={{
                          // Code component - preserve className for language detection
                          code({ children, className, ...props }) {
                            return (
                              <code className={className} {...props}>
                                {children}
                              </code>
                            )
                          },
                          // Code block - use SyntaxHighlighter wrapped in div
                          pre({ children }) {
                            // Extract code content and language from children
                            const codeElement = children as React.ReactElement
                            const className = codeElement?.props?.className || ''
                            const match = /language-(\w+)/.exec(className)
                            const codeContent = codeElement?.props?.children || ''

                            return (
                              <div className="not-prose my-6">
                                <SyntaxHighlighter
                                  style={oneDark}
                                  language={match ? match[1] : 'text'}
                                  PreTag="pre"
                                  customStyle={{
                                    margin: 0,
                                    padding: '1rem',
                                    borderRadius: '0.5rem',
                                    border: '1px solid hsl(217 33% 25%)',
                                    fontSize: '0.875rem',
                                    background: 'hsl(222 47% 9%)',
                                  }}
                                  codeTagProps={{
                                    style: {
                                      fontFamily: "'Fira Code', 'Consolas', monospace",
                                    },
                                  }}
                                >
                                  {String(codeContent).replace(/\n$/, '')}
                                </SyntaxHighlighter>
                              </div>
                            )
                          },
                          // Custom link rendering - resolve relative paths to GitHub URLs
                          a({ href, children, ...props }) {
                            let resolvedHref = href
                            if (isRelativePath(href)) {
                              resolvedHref = resolveGitHubUrl(
                                href!,
                                info.environment.path,
                                GITHUB_BLOB_BASE
                              )
                            }
                            return (
                              <a
                                href={resolvedHref}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-primary"
                                {...props}
                              >
                                {children}
                              </a>
                            )
                          },
                          // Custom heading rendering
                          h1({ children, ...props }) {
                            return (
                              <h1
                                className="mb-4 mt-8 border-b border-border/50 pb-2 text-2xl font-bold text-foreground first:mt-0"
                                {...props}
                              >
                                {children}
                              </h1>
                            )
                          },
                          h2({ children, ...props }) {
                            return (
                              <h2
                                className="mb-3 mt-8 border-b border-border/30 pb-2 text-xl font-semibold text-foreground"
                                {...props}
                              >
                                {children}
                              </h2>
                            )
                          },
                          h3({ children, ...props }) {
                            return (
                              <h3
                                className="mb-2 mt-6 text-lg font-semibold text-foreground"
                                {...props}
                              >
                                {children}
                              </h3>
                            )
                          },
                          // Custom paragraph
                          p({ children, ...props }) {
                            return (
                              <p
                                className="my-4 leading-7 text-muted-foreground"
                                {...props}
                              >
                                {children}
                              </p>
                            )
                          },
                          // Custom list rendering
                          ul({ children, ...props }) {
                            return (
                              <ul
                                className="my-4 ml-6 list-disc space-y-2 text-muted-foreground"
                                {...props}
                              >
                                {children}
                              </ul>
                            )
                          },
                          ol({ children, ...props }) {
                            return (
                              <ol
                                className="my-4 ml-6 list-decimal space-y-2 text-muted-foreground"
                                {...props}
                              >
                                {children}
                              </ol>
                            )
                          },
                          li({ children, ...props }) {
                            return (
                              <li className="leading-7" {...props}>
                                {children}
                              </li>
                            )
                          },
                          // Custom blockquote
                          blockquote({ children, ...props }) {
                            return (
                              <blockquote
                                className="my-6 border-l-4 border-primary bg-muted/30 py-2 pl-4 pr-4 italic text-muted-foreground"
                                {...props}
                              >
                                {children}
                              </blockquote>
                            )
                          },
                          // Custom table rendering
                          table({ children, ...props }) {
                            return (
                              <div className="my-6 overflow-x-auto">
                                <table
                                  className="w-full border-collapse text-sm"
                                  {...props}
                                >
                                  {children}
                                </table>
                              </div>
                            )
                          },
                          thead({ children, ...props }) {
                            return (
                              <thead
                                className="border-b-2 border-border/50 bg-muted/30"
                                {...props}
                              >
                                {children}
                              </thead>
                            )
                          },
                          th({ children, ...props }) {
                            return (
                              <th
                                className="px-4 py-3 text-left font-semibold text-foreground"
                                {...props}
                              >
                                {children}
                              </th>
                            )
                          },
                          td({ children, ...props }) {
                            return (
                              <td
                                className="border-b border-border/30 px-4 py-3 text-muted-foreground"
                                {...props}
                              >
                                {children}
                              </td>
                            )
                          },
                          // Custom hr
                          hr({ ...props }) {
                            return (
                              <hr
                                className="my-8 border-t border-border/50"
                                {...props}
                              />
                            )
                          },
                          // Custom strong/bold
                          strong({ children, ...props }) {
                            return (
                              <strong
                                className="font-semibold text-foreground"
                                {...props}
                              >
                                {children}
                              </strong>
                            )
                          },
                          // Custom em/italic
                          em({ children, ...props }) {
                            return (
                              <em className="italic text-muted-foreground" {...props}>
                                {children}
                              </em>
                            )
                          },
                          // Image - resolve relative paths to GitHub raw URLs
                          img({ src, alt, ...props }) {
                            let resolvedSrc = src
                            if (isRelativePath(src)) {
                              resolvedSrc = resolveGitHubUrl(
                                src!,
                                info.environment.path,
                                GITHUB_RAW_BASE
                              )
                            }
                            return (
                              <img
                                referrerPolicy="no-referrer"
                                src={resolvedSrc}
                                alt={alt || ''}
                                className="my-6 max-w-full rounded-lg border border-border/50"
                                {...props}
                              />
                            )
                          },
                        }}
                      >
                        {info.readme}
                      </ReactMarkdown>
                    </article>
                ) : (
                  <p className="font-mono text-sm text-muted-foreground">
                    No README available
                  </p>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="compose">
            <Card className="border-border/50 bg-card/80 backdrop-blur">
              <CardHeader className="border-b border-border/50">
                <CardTitle className="flex items-center gap-2 font-mono text-sm uppercase tracking-wider">
                  <FileCode className="h-4 w-4 text-primary" />
                  docker-compose.yml
                </CardTitle>
              </CardHeader>
              <CardContent className="pt-4">
                {info.compose_file ? (
                  <SyntaxHighlighter
                      style={oneDark}
                      language="yaml"
                      showLineNumbers
                      customStyle={{
                        background: 'hsl(222 47% 9%)',
                        borderRadius: '0.5rem',
                        border: '1px solid hsl(217 33% 25%)',
                        fontSize: '0.875rem',
                        margin: 0,
                      }}
                      codeTagProps={{
                        style: {
                          fontFamily: "'Fira Code', monospace",
                        },
                      }}
                    >
                      {info.compose_file}
                    </SyntaxHighlighter>
                ) : (
                  <p className="font-mono text-sm text-muted-foreground">
                    Compose file not available. Start the environment to download
                    it.
                  </p>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
