"use client"

import { useState } from "react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { BreadcrumbProvider } from "@/contexts/breadcrumb"
import { BreadcrumbSidebarLayout } from "@/components/breadcrumb-sidebar-layout"

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient())
  return (
    <QueryClientProvider client={queryClient}>
      <BreadcrumbProvider>
        <BreadcrumbSidebarLayout>{children}</BreadcrumbSidebarLayout>
      </BreadcrumbProvider>
    </QueryClientProvider>
  )
}
