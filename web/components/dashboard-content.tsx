"use client"

import * as React from "react"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  UserGroupIcon,
  AnalyticsUpIcon,
  ArrowUp01Icon,
  ArrowDown01Icon,
  Folder01Icon,
  Calendar03Icon,
  Tick02Icon,
  Clock01Icon,
} from "@hugeicons/core-free-icons"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Progress } from "@/components/ui/progress"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

const statsData = [
  {
    title: "Total de Usuarios",
    value: "2,847",
    change: "+12.5%",
    trend: "up" as const,
    icon: UserGroupIcon,
  },
  {
    title: "Receita Mensal",
    value: "R$ 54.320",
    change: "+8.2%",
    trend: "up" as const,
    icon: AnalyticsUpIcon,
  },
  {
    title: "Projetos Ativos",
    value: "24",
    change: "-2.4%",
    trend: "down" as const,
    icon: Folder01Icon,
  },
  {
    title: "Tarefas Concluidas",
    value: "1,284",
    change: "+18.7%",
    trend: "up" as const,
    icon: Tick02Icon,
  },
]

const recentActivity = [
  {
    user: "Sofia Andrade",
    avatar: "/avatars/sofia.jpg",
    initials: "SA",
    action: "criou um novo projeto",
    target: "Design System v2",
    time: "ha 5 minutos",
  },
  {
    user: "Carlos Mendes",
    avatar: "/avatars/carlos.jpg",
    initials: "CM",
    action: "comentou em",
    target: "Revisao de codigo #127",
    time: "ha 12 minutos",
  },
  {
    user: "Ana Silva",
    avatar: "/avatars/ana.jpg",
    initials: "AS",
    action: "concluiu a tarefa",
    target: "Implementar sidebar",
    time: "ha 1 hora",
  },
  {
    user: "Pedro Costa",
    avatar: "/avatars/pedro.jpg",
    initials: "PC",
    action: "adicionou arquivo em",
    target: "Documentacao",
    time: "ha 2 horas",
  },
]

const upcomingTasks = [
  {
    title: "Reuniao de planejamento",
    time: "14:00",
    date: "Hoje",
    priority: "high" as const,
  },
  {
    title: "Revisao de design",
    time: "10:00",
    date: "Amanha",
    priority: "medium" as const,
  },
  {
    title: "Deploy de producao",
    time: "16:00",
    date: "Amanha",
    priority: "high" as const,
  },
  {
    title: "Atualizacao de dependencias",
    time: "09:00",
    date: "Sexta",
    priority: "low" as const,
  },
]

const projectProgress = [
  { name: "Design System", progress: 85 },
  { name: "Aplicacao Web", progress: 62 },
  { name: "API Backend", progress: 45 },
  { name: "Documentacao", progress: 90 },
]

export function DashboardContent() {
  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Stats Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statsData.map((stat) => (
          <Card key={stat.title}>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div className="flex flex-col gap-1">
                  <span className="text-xs text-muted-foreground">
                    {stat.title}
                  </span>
                  <span className="text-2xl font-bold">{stat.value}</span>
                  <div className="flex items-center gap-1">
                    <HugeiconsIcon
                      icon={stat.trend === "up" ? ArrowUp01Icon : ArrowDown01Icon}
                      strokeWidth={2}
                      className={stat.trend === "up" ? "text-emerald-500" : "text-rose-500"}
                    />
                    <span
                      className={`text-xs ${
                        stat.trend === "up" ? "text-emerald-500" : "text-rose-500"
                      }`}
                    >
                      {stat.change}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      vs. mes anterior
                    </span>
                  </div>
                </div>
                <div className="flex size-12 items-center justify-center rounded-lg bg-primary/10">
                  <HugeiconsIcon
                    icon={stat.icon}
                    strokeWidth={2}
                    className="size-6 text-primary"
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Main Content Grid */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Recent Activity */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Atividade Recente</CardTitle>
            <CardDescription>
              Ultimas acoes realizadas pela equipe
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-4">
              {recentActivity.map((activity, index) => (
                <div
                  key={index}
                  className="flex items-center gap-4 rounded-lg border p-3"
                >
                  <Avatar size="sm">
                    <AvatarImage src={activity.avatar} alt={activity.user} />
                    <AvatarFallback>{activity.initials}</AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm">
                      <span className="font-medium">{activity.user}</span>{" "}
                      <span className="text-muted-foreground">{activity.action}</span>{" "}
                      <span className="font-medium text-primary">{activity.target}</span>
                    </p>
                    <p className="text-xs text-muted-foreground">{activity.time}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Upcoming Tasks */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Proximas Tarefas</CardTitle>
                <CardDescription>Agenda da semana</CardDescription>
              </div>
              <Button variant="ghost" size="icon-sm">
                <HugeiconsIcon icon={Calendar03Icon} strokeWidth={2} />
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-3">
              {upcomingTasks.map((task, index) => (
                <div
                  key={index}
                  className="flex items-center gap-3 rounded-lg border p-3"
                >
                  <div className="flex size-8 items-center justify-center rounded-lg bg-muted">
                    <HugeiconsIcon icon={Clock01Icon} strokeWidth={2} className="size-4" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{task.title}</p>
                    <p className="text-xs text-muted-foreground">
                      {task.date} as {task.time}
                    </p>
                  </div>
                  <Badge
                    variant={
                      task.priority === "high"
                        ? "default"
                        : task.priority === "medium"
                        ? "secondary"
                        : "outline"
                    }
                    className="text-xs"
                  >
                    {task.priority === "high"
                      ? "Alta"
                      : task.priority === "medium"
                      ? "Media"
                      : "Baixa"}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Project Progress */}
      <Card>
        <CardHeader>
          <CardTitle>Progresso dos Projetos</CardTitle>
          <CardDescription>
            Acompanhamento do desenvolvimento dos projetos ativos
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {projectProgress.map((project) => (
              <div key={project.name} className="flex flex-col gap-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">{project.name}</span>
                  <span className="text-sm text-muted-foreground">
                    {project.progress}%
                  </span>
                </div>
                <Progress value={project.progress} className="h-2" />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
