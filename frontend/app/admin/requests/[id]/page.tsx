"use client"

import { useEffect, useState, use } from "react"
import { useRouter } from "next/navigation"
import { useStore } from "@/store/use-store"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { ArrowLeft, User as UserIcon } from "lucide-react"
import Link from "next/link"
import { format, parseISO } from "date-fns"
import { ru } from "date-fns/locale"

export default function AdminRequestPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const router = useRouter()
  const { requests, fetchAllRequests, mentors, fetchMentors, assignMentor, isLoading } = useStore()
  
  const [selectedMentor, setSelectedMentor] = useState<string>("")

  useEffect(() => {
    fetchAllRequests()
    fetchMentors()
  }, [fetchAllRequests, fetchMentors])

  const request = (requests || []).find(r => r.id === id)

  if (!request) return <div>Загрузка...</div>

  const handleAssign = async () => {
    if (!selectedMentor) return
    await assignMentor(request.id, selectedMentor)
    router.push('/admin')
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/admin">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold text-slate-900">Детали заявки</h1>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <div className="flex justify-between items-start">
              <div>
                <CardTitle>{request.topic}</CardTitle>
                <CardDescription>
                  Создана: {request.createdAt && format(parseISO(request.createdAt), 'd MMMM yyyy', { locale: ru })}
                </CardDescription>
              </div>
              <Badge variant={
                request.status === 'pending' ? 'secondary' : 
                request.status === 'approved' ? 'success' : 'destructive'
              }>
                {request.status === 'pending' && 'На рассмотрении'}
                {request.status === 'approved' && 'Назначена'}
                {request.status === 'rejected' && 'Отклонена'}
              </Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <h3 className="text-sm font-medium text-slate-500 mb-1">Сотрудник</h3>
              <div className="flex items-center gap-2">
                <div className="h-8 w-8 bg-slate-100 rounded-full flex items-center justify-center">
                  <UserIcon className="h-4 w-4 text-slate-500" />
                </div>
                <div>
                  <p className="font-medium">{request.user?.name}</p>
                  <p className="text-xs text-slate-500">{request.user?.jobTitle}</p>
                </div>
              </div>
            </div>
            
            <div>
              <h3 className="text-sm font-medium text-slate-500 mb-1">Описание и мотивация</h3>
              <div className="p-4 bg-[#F2F3F7] rounded-xl text-sm text-slate-700">
                {request.description}
              </div>
            </div>
          </CardContent>
        </Card>

        {request.status === 'pending' && (
          <Card className="border-slate-300 shadow-md">
            <CardHeader>
              <CardTitle>Назначить наставника</CardTitle>
              <CardDescription>Выберите подходящего наставника из списка</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Наставник</Label>
                <Select onValueChange={setSelectedMentor} value={selectedMentor}>
                  <SelectTrigger>
                    <SelectValue placeholder="Выберите наставника" />
                  </SelectTrigger>
                  <SelectContent>
                    {(mentors || []).map((mentor) => (
                      <SelectItem key={mentor.id} value={mentor.id}>
                        <div className="flex flex-col items-start py-1">
                          <span className="font-medium">{mentor.name}</span>
                          <span className="text-xs text-slate-500">
                            {mentor.jobTitle} • Загрузка: {mentor.workload}/5
                          </span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              
              {selectedMentor && (
                <div className="p-4 bg-blue-50 text-blue-900 rounded-xl text-sm">
                  <p className="font-medium">Выбрано:</p>
                  <p>{mentors.find(m => m.id === selectedMentor)?.name}</p>
                  <p className="text-xs mt-1 opacity-80">После назначения сотруднику придет уведомление, и создастся процесс обучения.</p>
                </div>
              )}
            </CardContent>
            <CardFooter className="flex justify-end gap-3 bg-[#F2F3F7] border-t p-4">
              <Button variant="outline" asChild>
                <Link href="/admin">Отмена</Link>
              </Button>
              <Button onClick={handleAssign} disabled={!selectedMentor || isLoading}>
                {isLoading ? "Назначение..." : "Назначить"}
              </Button>
            </CardFooter>
          </Card>
        )}
      </div>
    </div>
  )
}

