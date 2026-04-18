"use client"

import * as React from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  Home01Icon,
  Folder01Icon,
  Settings01Icon,
  UserIcon,
  UserGroupIcon,
  Search01Icon,
  MailOpen01Icon,
  Calendar03Icon,
  Notification01Icon,
  ArrowRight01Icon,
  PlusSignIcon,
  MoreHorizontalIcon,
  Logout02Icon,
  CreditCardIcon,
  HelpCircleIcon,
  FavouriteIcon,
  AnalyticsUpIcon,
  LayersIcon,
  GridIcon,
  CommandIcon,
} from "@hugeicons/core-free-icons"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInput,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
  SidebarSeparator,
  useSidebar,
} from "@/components/ui/sidebar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"

// Navigation data structure
const navigationData = {
  user: {
    name: "Sofia Andrade",
    email: "sofia@radix-lyra.com",
    avatar: "/avatars/sofia.jpg",
    initials: "SA",
  },
  teams: [
    {
      name: "Radix Lyra",
      logo: GridIcon,
      plan: "Enterprise",
    },
    {
      name: "Studio Alpha",
      logo: LayersIcon,
      plan: "Pro",
    },
    {
      name: "Dev Team",
      logo: CommandIcon,
      plan: "Free",
    },
  ],
  mainNav: [
    {
      title: "Inicio",
      url: "/",
      icon: Home01Icon,
      isActive: true,
    },
    {
      title: "Pesquisa",
      url: "/search",
      icon: Search01Icon,
    },
    {
      title: "Notificacoes",
      url: "/notifications",
      icon: Notification01Icon,
      badge: "12",
    },
    {
      title: "Mensagens",
      url: "/messages",
      icon: MailOpen01Icon,
      badge: "3",
    },
  ],
  projects: [
    {
      title: "Design System",
      url: "/projects/design-system",
      icon: LayersIcon,
      isActive: true,
      items: [
        { title: "Componentes", url: "/projects/design-system/components" },
        { title: "Tokens", url: "/projects/design-system/tokens" },
        { title: "Documentacao", url: "/projects/design-system/docs" },
      ],
    },
    {
      title: "Aplicacao Web",
      url: "/projects/web-app",
      icon: GridIcon,
      items: [
        { title: "Dashboard", url: "/projects/web-app/dashboard" },
        { title: "Configuracoes", url: "/projects/web-app/settings" },
      ],
    },
    {
      title: "Analytics",
      url: "/projects/analytics",
      icon: AnalyticsUpIcon,
      items: [
        { title: "Relatorios", url: "/projects/analytics/reports" },
        { title: "Metricas", url: "/projects/analytics/metrics" },
      ],
    },
  ],
  favorites: [
    {
      title: "Componentes UI",
      url: "/favorites/ui-components",
      icon: FavouriteIcon,
    },
    {
      title: "Biblioteca de Icones",
      url: "/favorites/icons",
      icon: FavouriteIcon,
    },
  ],
  secondaryNav: [
    {
      title: "Calendario",
      url: "/calendar",
      icon: Calendar03Icon,
    },
    {
      title: "Arquivos",
      url: "/files",
      icon: Folder01Icon,
    },
    {
      title: "Configuracoes",
      url: "/settings",
      icon: Settings01Icon,
    },
  ],
}

function TeamSwitcher() {
  const [activeTeam, setActiveTeam] = React.useState(navigationData.teams[0])
  const { isMobile } = useSidebar()

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <HugeiconsIcon icon={activeTeam.logo} strokeWidth={2} className="size-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">{activeTeam.name}</span>
                <span className="truncate text-xs text-muted-foreground">{activeTeam.plan}</span>
              </div>
              <HugeiconsIcon icon={ArrowRight01Icon} strokeWidth={2} className="ml-auto rotate-90" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56"
            align="start"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel>Equipes</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {navigationData.teams.map((team) => (
              <DropdownMenuItem
                key={team.name}
                onClick={() => setActiveTeam(team)}
                className="gap-2 p-2"
              >
                <div className="flex size-6 items-center justify-center rounded-sm border bg-background">
                  <HugeiconsIcon icon={team.logo} strokeWidth={2} className="size-4" />
                </div>
                <span className="flex-1">{team.name}</span>
                <Badge variant="secondary" className="text-xs">
                  {team.plan}
                </Badge>
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className="gap-2 p-2">
              <div className="flex size-6 items-center justify-center rounded-sm border bg-background">
                <HugeiconsIcon icon={PlusSignIcon} strokeWidth={2} className="size-4" />
              </div>
              <span className="text-muted-foreground">Adicionar equipe</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

const clientesNav = [
  { title: "Todos os clientes", url: "/clientes" },
  { title: "Novo cliente", url: "/clientes/novo" },
]

function NavClientes() {
  const pathname = usePathname()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Clientes</SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          <Collapsible defaultOpen>
            <SidebarMenuItem>
              <SidebarMenuButton asChild isActive={pathname.startsWith("/clientes")} tooltip="Clientes">
                <Link href="/clientes">
                  <HugeiconsIcon icon={UserGroupIcon} strokeWidth={2} />
                  <span>Clientes</span>
                </Link>
              </SidebarMenuButton>
              <CollapsibleTrigger asChild>
                <SidebarMenuAction className="data-open:rotate-90 transition-transform">
                  <HugeiconsIcon icon={ArrowRight01Icon} strokeWidth={2} />
                  <span className="sr-only">Expandir</span>
                </SidebarMenuAction>
              </CollapsibleTrigger>
              <CollapsibleContent>
                <SidebarMenuSub>
                  {clientesNav.map((item) => (
                    <SidebarMenuSubItem key={item.title}>
                      <SidebarMenuSubButton asChild isActive={pathname === item.url}>
                        <Link href={item.url}>
                          <span>{item.title}</span>
                        </Link>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                  ))}
                </SidebarMenuSub>
              </CollapsibleContent>
            </SidebarMenuItem>
          </Collapsible>
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

function NavMain() {
  const pathname = usePathname()

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Navegacao</SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          {navigationData.mainNav.map((item) => {
            const isActive = item.url === "/" ? pathname === "/" : pathname.startsWith(item.url)
            return (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton asChild isActive={isActive} tooltip={item.title}>
                  <Link href={item.url}>
                    <HugeiconsIcon icon={item.icon} strokeWidth={2} />
                    <span>{item.title}</span>
                  </Link>
                </SidebarMenuButton>
                {item.badge && (
                  <SidebarMenuBadge>{item.badge}</SidebarMenuBadge>
                )}
              </SidebarMenuItem>
            )
          })}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

function NavProjects() {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>Projetos</SidebarGroupLabel>
      <SidebarGroupAction title="Adicionar Projeto">
        <HugeiconsIcon icon={PlusSignIcon} strokeWidth={2} />
        <span className="sr-only">Adicionar Projeto</span>
      </SidebarGroupAction>
      <SidebarGroupContent>
        <SidebarMenu>
          {navigationData.projects.map((project) => (
            <Collapsible key={project.title} asChild defaultOpen={project.isActive}>
              <SidebarMenuItem>
                <SidebarMenuButton asChild tooltip={project.title}>
                  <Link href={project.url}>
                    <HugeiconsIcon icon={project.icon} strokeWidth={2} />
                    <span>{project.title}</span>
                  </Link>
                </SidebarMenuButton>
                {project.items?.length ? (
                  <>
                    <CollapsibleTrigger asChild>
                      <SidebarMenuAction className="data-open:rotate-90 transition-transform">
                        <HugeiconsIcon icon={ArrowRight01Icon} strokeWidth={2} />
                        <span className="sr-only">Expandir</span>
                      </SidebarMenuAction>
                    </CollapsibleTrigger>
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        {project.items.map((subItem) => (
                          <SidebarMenuSubItem key={subItem.title}>
                            <SidebarMenuSubButton asChild>
                              <Link href={subItem.url}>
                                <span>{subItem.title}</span>
                              </Link>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        ))}
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  </>
                ) : null}
              </SidebarMenuItem>
            </Collapsible>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

function NavFavorites() {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>Favoritos</SidebarGroupLabel>
      <SidebarGroupContent>
        <SidebarMenu>
          {navigationData.favorites.map((item) => (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild tooltip={item.title}>
                <Link href={item.url}>
                  <HugeiconsIcon icon={item.icon} strokeWidth={2} className="text-amber-500" />
                  <span>{item.title}</span>
                </Link>
              </SidebarMenuButton>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <SidebarMenuAction showOnHover>
                    <HugeiconsIcon icon={MoreHorizontalIcon} strokeWidth={2} />
                    <span className="sr-only">Mais opcoes</span>
                  </SidebarMenuAction>
                </DropdownMenuTrigger>
                <DropdownMenuContent side="right" align="start">
                  <DropdownMenuItem>
                    <span>Remover dos favoritos</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <span>Copiar link</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

function NavSecondary() {
  return (
    <SidebarGroup className="mt-auto">
      <SidebarGroupContent>
        <SidebarMenu>
          {navigationData.secondaryNav.map((item) => (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild tooltip={item.title}>
                <Link href={item.url}>
                  <HugeiconsIcon icon={item.icon} strokeWidth={2} />
                  <span>{item.title}</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

function NavUser() {
  const { isMobile } = useSidebar()
  const { user } = navigationData

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar size="sm">
                <AvatarImage src={user.avatar} alt={user.name} />
                <AvatarFallback>{user.initials}</AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">{user.name}</span>
                <span className="truncate text-xs text-muted-foreground">{user.email}</span>
              </div>
              <HugeiconsIcon icon={ArrowRight01Icon} strokeWidth={2} className="ml-auto rotate-90" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5">
                <Avatar size="sm">
                  <AvatarImage src={user.avatar} alt={user.name} />
                  <AvatarFallback>{user.initials}</AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">{user.name}</span>
                  <span className="truncate text-xs text-muted-foreground">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <HugeiconsIcon icon={UserIcon} strokeWidth={2} />
                <span>Perfil</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <HugeiconsIcon icon={CreditCardIcon} strokeWidth={2} />
                <span>Cobranca</span>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <HugeiconsIcon icon={Settings01Icon} strokeWidth={2} />
                <span>Configuracoes</span>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <HugeiconsIcon icon={HelpCircleIcon} strokeWidth={2} />
                <span>Ajuda</span>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem variant="destructive">
              <HugeiconsIcon icon={Logout02Icon} strokeWidth={2} />
              <span>Sair</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher />
        <div className="px-2">
          <SidebarInput placeholder="Pesquisar..." />
        </div>
      </SidebarHeader>
      <SidebarSeparator />
      <SidebarContent>
        <NavClientes />
        <SidebarSeparator />
        <NavMain />
        <SidebarSeparator />
        <NavProjects />
        <SidebarSeparator />
        <NavFavorites />
        <NavSecondary />
      </SidebarContent>
      <SidebarSeparator />
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
