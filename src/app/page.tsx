import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";

export default function Home() {
  return (
    <main className="mx-auto flex min-h-full w-full max-w-md flex-col justify-center gap-6 px-6 py-16">
      <p className="text-sm tracking-wide text-neutral-500">SyncLink</p>
      <h1 className="text-3xl font-medium tracking-tight">Your links. One page.</h1>
      <p className="text-neutral-600">Text-first public pages. Edit them from the dashboard.</p>
      <Link href="/dashboard" className={buttonVariants()}>Open dashboard</Link>
    </main>
  );
}
