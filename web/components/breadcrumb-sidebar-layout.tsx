"use client"

import { useBreadcrumbs } from "@/contexts/breadcrumb"
import { SidebarLayout } from "@/components/sidebar-layout"

export function BreadcrumbSidebarLayout({ children }: { children: React.ReactNode }) {
  const breadcrumbs = useBreadcrumbs()
  return <SidebarLayout breadcrumbs={breadcrumbs}>{children}</SidebarLayout>
}
