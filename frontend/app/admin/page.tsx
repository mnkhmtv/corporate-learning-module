"use client"

import { useEffect } from "react"
import Link from "next/link"
import { useStore } from "@/store/use-store"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { format } from "date-fns"
import { ru } from "date-fns/locale"

export default function AdminDashboardPage() {
  const { requests, fetchAllRequests, isLoading } = useStore()

  useEffect(() => {
    fetchAllRequests()
  }, [fetchAllRequests])

  const pendingRequests = requests.filter(r => r.status === 'pending')
  const otherRequests = requests.filter(r => r.status !== 'pending')

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-slate-900">Заявки на обучение</h1>
      </div>

      {/* New Requests */}
      <section>
        <h2 className="text-lg font-semibold mb-4">Новые заявки ({pendingRequests.length})</h2>
        <div className="bg-white rounded-md border border-slate-200 overflow-hidden">
          <table className="w-full text-sm text-left">
            <thead className="bg-slate-50 text-slate-500 border-b border-slate-200">
              <tr>
                <th className="px-4 py-3 font-medium">Сотрудник</th>
                <th className="px-4 py-3 font-medium">Тема</th>
                <th className="px-4 py-3 font-medium">Дата</th>
                <th className="px-4 py-3 font-medium">Статус</th>
                <th className="px-4 py-3 font-medium text-right">Действие</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {pendingRequests.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-slate-500">
                    Нет новых заявок
                  </td>
                </tr>
              ) : (
                pendingRequests.map((request) => (
                  <tr key={request.id} className="hover:bg-slate-50/50">
                    <td className="px-4 py-3">
                      <div className="font-medium">User ID: {request.userId}</div>
                      {/* Ideally we would join with user table, but mock service simplifies this */}
                    </td>
                    <td className="px-4 py-3">{request.topic}</td>
                    <td className="px-4 py-3 text-slate-500">
                      {format(new Date(request.createdAt), 'd MMM yyyy', { locale: ru })}
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant="secondary">Новая</Badge>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <Button size="sm" asChild>
                        <Link href={`/admin/requests/${request.id}`}>
                          Открыть
                        </Link>
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      {/* History */}
      <section>
        <h2 className="text-lg font-semibold mb-4">История заявок</h2>
        <div className="bg-white rounded-md border border-slate-200 overflow-hidden">
          <table className="w-full text-sm text-left">
            <thead className="bg-slate-50 text-slate-500 border-b border-slate-200">
              <tr>
                <th className="px-4 py-3 font-medium">Сотрудник</th>
                <th className="px-4 py-3 font-medium">Тема</th>
                <th className="px-4 py-3 font-medium">Дата</th>
                <th className="px-4 py-3 font-medium">Статус</th>
                <th className="px-4 py-3 font-medium text-right">Действие</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
               {otherRequests.map((request) => (
                  <tr key={request.id} className="hover:bg-slate-50/50">
                    <td className="px-4 py-3">
                       <div className="font-medium">User ID: {request.userId}</div>
                    </td>
                    <td className="px-4 py-3">{request.topic}</td>
                    <td className="px-4 py-3 text-slate-500">
                      {format(new Date(request.createdAt), 'd MMM yyyy', { locale: ru })}
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant={request.status === 'approved' ? 'success' : 'destructive'}>
                        {request.status === 'approved' ? 'Назначена' : 'Отклонена'}
                      </Badge>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <Button variant="ghost" size="sm" asChild>
                        <Link href={`/admin/requests/${request.id}`}>
                          Детали
                        </Link>
                      </Button>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}

