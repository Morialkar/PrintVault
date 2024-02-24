'use client';
import { FuseProvider } from '../components/FuseProvider';
import './global.css';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <FuseProvider>
        <body>{children}</body>
      </FuseProvider>
    </html>
  );
}
