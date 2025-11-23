"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { useStore } from "@/store/use-store"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft } from "lucide-react"
import Link from "next/link"

export default function CreateRequestPage() {
  const router = useRouter()
  const createRequest = useStore((state) => state.createRequest)
  const isLoading = useStore((state) => state.isLoading)
  
  const [topic, setTopic] = useState("")
  const [description, setDescription] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await createRequest({ topic, description })
    router.push('/dashboard')
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-6">
        <Button variant="ghost" asChild className="mb-2 pl-0 hover:pl-2 transition-all">
          <Link href="/dashboard">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Назад
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Создание заявки</h1>
        <p className="text-slate-500">Опишите тему, которую вы хотите изучить</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Детали заявки</CardTitle>
          <CardDescription>
            Администратор рассмотрит вашу заявку и подберет подходящего наставника в течение 1-2 дней.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="topic">Тема обучения</Label>
              <Input
                id="topic"
                placeholder="Например: Продвинутые паттерны React"
                value={topic}
                onChange={(e) => setTopic(e.target.value)}
                required
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="description">Описание и мотивация</Label>
              <Textarea
                id="description"
                placeholder="Почему вам важно изучить эту тему? Какие задачи это поможет решить?"
                className="min-h-[150px]"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                required
              />
            </div>

            <div className="flex justify-end gap-4">
              <Button variant="outline" type="button" asChild>
                <Link href="/dashboard">Отмена</Link>
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Отправка..." : "Отправить заявку"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

