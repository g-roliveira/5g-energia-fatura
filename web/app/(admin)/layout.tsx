"use client"

import { useState } from "react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { SidebarLayout } from "@/components/sidebar-layout"

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient())
  return (
    <QueryClientProvider client={queryClient}>
      <SidebarLayout>{children}</SidebarLayout>
    </QueryClientProvider>
  )
}
