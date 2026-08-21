"use client";

import { FormEvent, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { synclink } from "@/lib/api";

export default function ResetPage() {
  const [email, setEmail] = useState("");
  const [token, setToken] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    try {
      await synclink.resetPassword(email, token, newPassword);
      setNotice("Password updated. You can sign in.");
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Reset failed");
    }
  }

  return (
    <main className="flex min-h-full items-center justify-center bg-[#faf9f7] px-6 py-16">
      <Card className="w-full max-w-md border-neutral-200/80 bg-white shadow-none">
        <CardHeader>
          <p className="text-xs tracking-[0.28em] text-neutral-400">STATION 03 — RESET</p>
          <CardTitle className="font-medium">New password</CardTitle>
          <CardDescription>Email, demo token, new password (min 8).</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={onSubmit}>
            <div className="space-y-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required /></div>
            <div className="space-y-2"><Label htmlFor="token">Token</Label><Input id="token" value={token} onChange={(event) => setToken(event.target.value)} required /></div>
            <div className="space-y-2"><Label htmlFor="password">New password</Label><Input id="password" type="password" minLength={8} value={newPassword} onChange={(event) => setNewPassword(event.target.value)} required /></div>
            {error ? <p className="text-sm text-red-600">{error}</p> : null}
            {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}
            <Button type="submit" className="w-full">Reset password</Button>
            <Button asChild variant="ghost" className="w-full"><Link href="/dashboard">Back to studio</Link></Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
