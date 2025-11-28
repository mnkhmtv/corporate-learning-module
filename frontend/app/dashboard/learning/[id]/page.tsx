"use client"

import { useEffect, useState, use } from "react"
import { useRouter } from "next/navigation"
import { useStore } from "@/store/use-store"
import { LearningPlanItem } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Textarea } from "@/components/ui/textarea"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogFooter, 
  DialogHeader, 
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog"
import { CheckCircle2, Mail, MessageCircle, ArrowLeft, Star, Plus, Trash2 } from "lucide-react"
import { cn } from "@/lib/utils"
import Link from "next/link"

export default function LearningPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const router = useRouter()
  const { fetchLearning, updateLearningPlan, completeLearning, updateNotes } = useStore()
  const [learning, setLearning] = useState<any>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [newPlanItem, setNewPlanItem] = useState("")
  const [notes, setNotes] = useState("")
  const [isSavingNotes, setIsSavingNotes] = useState(false)
  
  // Feedback state
  const [isFeedbackOpen, setIsFeedbackOpen] = useState(false)
  const [rating, setRating] = useState(0)
  const [comment, setComment] = useState("")

  useEffect(() => {
    const load = async () => {
      const data = await fetchLearning(id)
      if (data) {
        setLearning(data)
        setNotes(data.notes || "")
      } else {
        router.push('/dashboard')
      }
      setIsLoading(false)
    }
    load()
  }, [id, fetchLearning, router])

  if (isLoading || !learning) return <div>Loading...</div>

  const handleAddPlanItem = async () => {
    if (!newPlanItem.trim()) return
    const newItem: LearningPlanItem = {
      id: Math.random().toString(36).substr(2, 9),
      text: newPlanItem,
      completed: false
    }
    const updatedPlan = [...(learning.plan || []), newItem]
    await updateLearningPlan(id, updatedPlan)
    setLearning({ ...learning, plan: updatedPlan })
    setNewPlanItem("")
  }

  const togglePlanItem = async (itemId: string) => {
    const updatedPlan = learning.plan.map((item: LearningPlanItem) => 
      item.id === itemId ? { ...item, completed: !item.completed } : item
    )
    await updateLearningPlan(id, updatedPlan)
    setLearning({ ...learning, plan: updatedPlan })
  }

  const deletePlanItem = async (itemId: string) => {
    const updatedPlan = learning.plan.filter((item: LearningPlanItem) => item.id !== itemId)
    await updateLearningPlan(id, updatedPlan)
    setLearning({ ...learning, plan: updatedPlan })
  }

  const handleSaveNotes = async () => {
    setIsSavingNotes(true)
    await updateNotes(id, notes)
    setIsSavingNotes(false)
    // Optionally show a toast or confirmation message
  }

  const handleComplete = async () => {
    await completeLearning(id, { rating, comment })
    setIsFeedbackOpen(false)
    router.push('/dashboard')
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
         <Button variant="ghost" size="icon" asChild>
          <Link href="/dashboard">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-slate-900">{learning.topic}</h1>
          <div className="flex items-center gap-2 text-slate-500 text-sm">
            <Badge variant={learning.status === 'active' ? 'default' : 'secondary'}>
              {learning.status === 'active' ? 'В процессе' : 'Завершено'}
            </Badge>
            <span>•</span>
            <span>Наставник: {learning?.mentor.name}</span>
          </div>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Content - Plan */}
        <div className="md:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>План действий</CardTitle>
              <CardDescription>Составьте план обучения вместе с наставником</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-2">
                <Input 
                  placeholder="Добавить пункт плана..." 
                  value={newPlanItem}
                  onChange={(e) => setNewPlanItem(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && handleAddPlanItem()}
                />
                <Button onClick={handleAddPlanItem} size="icon" variant="secondary">
                  <Plus className="h-4 w-4" />
                </Button>
              </div>

              <div className="space-y-2">
                {learning.plan?.length === 0 && (
                  <p className="text-center text-slate-500 py-4">План пока пуст</p>
                )}
                {learning.plan?.map((item: LearningPlanItem) => (
                  <div key={item.id} className="flex items-center gap-3 p-3 border rounded-xl bg-[#F2F3F7]">
                    <button 
                      onClick={() => togglePlanItem(item.id)}
                      className={cn(
                        "flex-shrink-0 h-5 w-5 rounded border flex items-center justify-center transition-colors",
                        item.completed 
                          ? "bg-green-500 border-green-500 text-white" 
                          : "border-slate-300 bg-white"
                      )}
                    >
                      {item.completed && <CheckCircle2 className="h-3.5 w-3.5" />}
                    </button>
                    <span className={cn("flex-1 text-sm", item.completed && "text-slate-500 line-through")}>
                      {item.text}
                    </span>
                    <Button variant="ghost" size="icon" className="h-8 w-8 text-slate-400 hover:text-red-500" onClick={() => deletePlanItem(item.id)}>
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Мои заметки</CardTitle>
            </CardHeader>
            <CardContent>
              <Textarea 
                placeholder="Заметки по ходу обучения..." 
                className="min-h-[200px]"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
              />
              <Button 
                onClick={handleSaveNotes} 
                disabled={isSavingNotes}
                variant="secondary"
                className="mt-4"
              >
                {isSavingNotes ? 'Сохранение...' : 'Сохранить'}
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar - Mentor & Actions */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Наставник</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-600 font-bold">
                  {learning?.mentor.name.charAt(0)}
                </div>
                <div>
                  <p className="font-medium">{learning.mentor.name}</p>
                  <p className="text-xs text-slate-500">Ментор</p>
                </div>
              </div>
              
              <div className="space-y-2 pt-2 border-t text-sm">
                {learning.mentor.telegram && (
                   <div className="flex items-center gap-2 text-slate-600">
                    <MessageCircle className="h-4 w-4" />
                    <a href={`https://t.me/${learning.mentor.telegram.replace('@', '')}`} target="_blank" rel="noopener noreferrer" className="hover:underline">
                      {learning.mentor.telegram}
                    </a>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>

          <Dialog open={isFeedbackOpen} onOpenChange={setIsFeedbackOpen}>
            <DialogTrigger asChild>
              <Button className="w-full" size="lg" variant="default">
                Завершить обучение
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Завершение обучения</DialogTitle>
                <DialogDescription>
                  Оцените работу с наставником и оставьте отзыв, чтобы завершить обучение.
                </DialogDescription>
              </DialogHeader>
              
              <div className="space-y-4 py-4">
                <div className="space-y-2">
                  <Label>Оценка наставника</Label>
                  <div className="flex gap-2">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <button
                        key={star}
                        type="button"
                        onClick={() => setRating(star)}
                        className={cn(
                          "p-1 transition-colors",
                          rating >= star ? "text-yellow-400" : "text-slate-300 hover:text-yellow-200"
                        )}
                      >
                        <Star className="h-8 w-8 fill-current" />
                      </button>
                    ))}
                  </div>
                </div>
                
                <div className="space-y-2">
                  <Label>Ваш отзыв</Label>
                  <Textarea 
                    placeholder="Что понравилось? Что можно улучшить?"
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                  />
                </div>
              </div>

              <DialogFooter>
                <Button variant="outline" onClick={() => setIsFeedbackOpen(false)}>Отмена</Button>
                <Button onClick={handleComplete} disabled={rating === 0}>Отправить и завершить</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  )
}

