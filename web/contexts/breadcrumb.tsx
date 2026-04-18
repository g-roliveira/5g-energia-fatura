"use client"

import { createContext, useCallback, useContext, useEffect, useState } from "react"
import { usePathname } from "next/navigation"

const SEGMENT_LABELS: Record<string, string> = {
  clientes: "Clientes",
  novo: "Novo cliente",
  editar: "Editar",
  ucs: "UCs",
  faturas: "Faturas",
  sync: "Sincronização",
}

type BreadcrumbContextValue = {
  titles: Record<string, string>
  setTitle: (key: string, title: string) => void
}

const BreadcrumbContext = createContext<BreadcrumbContextValue>({
  titles: {},
  setTitle: () => {},
})

export function BreadcrumbProvider({ children }: { children: React.ReactNode }) {
  const [titles, setTitles] = useState<Record<string, string>>({})

  const setTitle = useCallback((key: string, title: string) => {
    setTitles((prev) => (prev[key] === title ? prev : { ...prev, [key]: title }))
  }, [])

  return (
    <BreadcrumbContext.Provider value={{ titles, setTitle }}>
      {children}
    </BreadcrumbContext.Provider>
  )
}

export function useBreadcrumbs() {
  const pathname = usePathname()
  const { titles } = useContext(BreadcrumbContext)

  const segments = pathname.split("/").filter(Boolean)
  const crumbs: { label: string; href?: string }[] = []
  let path = ""

  segments.forEach((seg, i) => {
    path += `/${seg}`
    const isLast = i === segments.length - 1
    const label = SEGMENT_LABELS[seg] ?? titles[seg] ?? seg
    crumbs.push(isLast ? { label } : { label, href: path })
  })

  return crumbs
}

export function useSetBreadcrumbTitle(key: string, title: string | null | undefined) {
  const { setTitle } = useContext(BreadcrumbContext)
  useEffect(() => {
    if (title) setTitle(key, title)
  }, [key, title, setTitle])
}
