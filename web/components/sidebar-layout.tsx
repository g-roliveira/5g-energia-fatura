"use client"

import * as React from "react"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  Search01Icon,
  Notification01Icon,
} from "@hugeicons/core-free-icons"

import {
  SidebarProvider,
  SidebarInset,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"
import { Separator } from "@/components/ui/separator"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Badge } from "@/components/ui/badge"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"

interface SidebarLayoutProps {
  children: React.ReactNode
  breadcrumbs?: {
    label: string
    href?: string
  }[]
  title?: string
}

export function SidebarLayout({
  children,
  breadcrumbs,
  title,
}: SidebarLayoutProps) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="sticky top-0 z-10 flex h-14 shrink-0 items-center gap-2 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="flex flex-1 items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            
            {breadcrumbs && breadcrumbs.length > 0 && (
              <Breadcrumb>
                <BreadcrumbList>
                  {breadcrumbs.map((crumb, index) => (
                    <React.Fragment key={crumb.label}>
                      <BreadcrumbItem className={index === 0 ? "hidden md:block" : ""}>
                        {index === breadcrumbs.length - 1 ? (
                          <BreadcrumbPage>{crumb.label}</BreadcrumbPage>
                        ) : (
                          <BreadcrumbLink href={crumb.href || "#"}>
                            {crumb.label}
                          </BreadcrumbLink>
                        )}
                      </BreadcrumbItem>
                      {index < breadcrumbs.length - 1 && (
                        <BreadcrumbSeparator className="hidden md:block" />
                      )}
                    </React.Fragment>
                  ))}
                </BreadcrumbList>
              </Breadcrumb>
            )}

            {title && !breadcrumbs && (
              <h1 className="text-sm font-medium">{title}</h1>
            )}

            <div className="ml-auto flex items-center gap-2">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="ghost" size="icon-sm">
                    <HugeiconsIcon icon={Search01Icon} strokeWidth={2} />
                    <span className="sr-only">Pesquisar</span>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Pesquisar</TooltipContent>
              </Tooltip>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon-sm" className="relative">
                    <HugeiconsIcon icon={Notification01Icon} strokeWidth={2} />
                    <Badge 
                      className="absolute -top-1 -right-1 h-4 min-w-4 px-1 text-[10px]"
                    >
                      5
                    </Badge>
                    <span className="sr-only">Notificacoes</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-80">
                  <div className="p-2">
                    <p className="text-sm font-medium">Notificacoes</p>
                    <p className="text-xs text-muted-foreground">
                      Voce tem 5 notificacoes nao lidas
                    </p>
                  </div>
                  <Separator />
                  <DropdownMenuItem className="flex flex-col items-start gap-1 p-3">
                    <span className="font-medium">Nova mensagem</span>
                    <span className="text-xs text-muted-foreground">
                      Sofia comentou no projeto Design System
                    </span>
                  </DropdownMenuItem>
                  <DropdownMenuItem className="flex flex-col items-start gap-1 p-3">
                    <span className="font-medium">Tarefa concluida</span>
                    <span className="text-xs text-muted-foreground">
                      A revisao de codigo foi aprovada
                    </span>
                  </DropdownMenuItem>
                  <Separator />
                  <DropdownMenuItem className="justify-center text-primary">
                    Ver todas as notificacoes
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </header>

        <main className="flex-1">
          {children}
        </main>
      </SidebarInset>
    </SidebarProvider>
  )
}
