import type { Metadata } from "next";
import { Fira_Code } from "next/font/google";
import "./globals.css";

const fira = Fira_Code({
  variable: "--font-fira",
  subsets: ["latin"],
  weight: ["400", "500", "600"],
});

export const metadata: Metadata = {
  title: "SyncLink",
  description: "One page. Every link.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className={`${fira.variable} h-full antialiased`}>
      <body className="min-h-full bg-[#faf9f7] font-sans text-neutral-950">{children}</body>
    </html>
  );
}
