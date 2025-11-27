"use client"

import { useEffect } from "react"
import Link from "next/link"
import { useStore } from "@/store/use-store"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Plus, BookOpen, Clock } from "lucide-react"
import { format, parseISO } from "date-fns"
import { ru } from "date-fns/locale"

export default function DashboardPage() {
  const { user, requests, learnings, fetchUserData } = useStore()

  useEffect(() => {
    fetchUserData()
  }, [fetchUserData])

  if (!user) return null

  const activeLearnings = (learnings || []).filter(l => l.status === 'active')
  const completedLearnings = (learnings || []).filter(l => l.status === 'completed')

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Личный кабинет</h1>
          <p className="text-slate-500">Добро пожаловать, {user.name}</p>
        </div>
        <Button asChild>
          <Link href="/dashboard/requests/create">
            <Plus className="mr-2 h-4 w-4" />
            Новая заявка
          </Link>
        </Button>
      </div>

      {/* Active Learnings */}
      <section>
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <BookOpen className="h-5 w-5" />
          Активные обучения
        </h2>
        {activeLearnings.length > 0 ? (
          <div className="grid gap-4 md:grid-cols-2">
            {activeLearnings.map((learning) => (
              <Card key={learning.id} className="hover:shadow-md transition-shadow">
                <CardHeader className="pb-2">
                  <div className="flex justify-between items-start">
                    <CardTitle className="text-xl">{learning.topic}</CardTitle>
                    <Badge variant="default">В процессе</Badge>
                  </div>
                  <CardDescription>Наставник: {learning?.mentorName}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-sm text-slate-500 mb-4">
                    Начало: {learning.startDate && format(parseISO(learning.startDate), 'd MMMM yyyy', { locale: ru })}
                  </div>
                  <Button className="w-full" asChild>
                    <Link href={`/dashboard/learning/${learning.id}`}>
                      Продолжить обучение
                    </Link>
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <Card className="bg-[#F2F3F7] border-dashed">
            <CardContent className="flex flex-col items-center justify-center py-8 text-center">
              <p className="text-slate-500 mb-4">У вас пока нет активных обучений</p>
              <Button variant="outline" asChild>
                <Link href="/dashboard/requests/create">Создать заявку</Link>
              </Button>
            </CardContent>
          </Card>
        )}
      </section>

      {/* History / Requests */}
      <section>
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <Clock className="h-5 w-5" />
          История заявок
        </h2>
        <div className="space-y-4">
          {(requests || []).map((request) => {
            const learning = learnings?.find(l => l.requestId === request.id);
            const isCompleted = learning?.status === 'completed';

            return (
              <Card key={request.id}>
                <CardContent className="flex items-center justify-between p-4">
                  <div>
                    <h3 className="font-medium">{request.topic}</h3>
                    <p className="text-sm text-slate-500">
                      Создана: {request.createdAt && format(parseISO(request.createdAt), 'd MMM yyyy', { locale: ru })}
                    </p>
                  </div>
                  <div className="flex items-center gap-4">
                    <Badge variant={
                      request.status === 'approved' ? (isCompleted ? 'outline' : 'success') : 
                      request.status === 'rejected' ? 'destructive' : 
                      'secondary'
                    }>
                      {request.status === 'pending' && 'На рассмотрении'}
                      {request.status === 'approved' && (isCompleted ? 'Закончено' : 'Назначено')}
                      {request.status === 'rejected' && 'Отклонено'}
                    </Badge>
                    {request.status === 'approved' && !isCompleted && (
                      <Button variant="ghost" size="sm" asChild>
                        <Link href={`/dashboard/learning/${learning?.id}`}>
                          Перейти
                        </Link>
                      </Button>
                    )}
                  </div>
                </CardContent>
              </Card>
            )
          })}
          {(!requests || requests.length === 0) && (
            <p className="text-slate-500 text-sm">История заявок пуста</p>
          )}
        </div>
      </section>

      {/* Completed History */}
      {completedLearnings.length > 0 && (
        <section>
           <h2 className="text-lg font-semibold mb-4">Завершенные обучения</h2>
           <div className="space-y-4">
             {completedLearnings.map((learning) => (
               <Card key={learning.id} className="bg-[#F2F3F7]">
                 <CardContent className="flex items-center justify-between p-4">
                   <div>
                     <h3 className="font-medium text-slate-700">{learning.topic}</h3>
                     <p className="text-sm text-slate-500">
                       Наставник: {learning?.mentorName}
                     </p>
                   </div>
                   <div className="flex items-center gap-2">
                      <div className="flex">
                        {[...Array(learning.feedback?.rating || 0)].map((_, i) => (
                          <span key={i} className="text-yellow-400">★</span>
                        ))}
                      </div>
                      <Badge variant="outline">Завершено</Badge>
                   </div>
                 </CardContent>
               </Card>
             ))}
           </div>
        </section>
      )}
    </div>
  )
}

