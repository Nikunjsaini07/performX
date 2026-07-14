'use client';

import React from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from 'recharts';

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{ value: number; payload: { player: string; team: string; flag: string } }>;
}

function CustomTooltip({ active, payload }: CustomTooltipProps) {
  if (!active || !payload || !payload.length) return null;
  const d = payload[0];
  return (
    <div className="bg-card border border-border rounded-lg px-4 py-3 shadow-xl">
      <p className="text-sm font-semibold text-foreground">
        {d.payload.flag} {d.payload.player}
      </p>
      <p className="text-xs text-muted-foreground">{d.payload.team}</p>
      <p className="text-sm font-bold text-primary mt-1 stat-number">
        {d.value.toFixed(1)} score
      </p>
    </div>
  );
}

interface Props {
  data: { player: string; team: string; value: number; flag: string }[];
}

export default function StatLeadersChart({ data }: Props) {
  if (!data || data.length === 0) {
    return <div className="p-8 text-center text-muted-foreground border border-border rounded-2xl">No trending data available.</div>;
  }

  return (
    <div className="archive-card p-6">
      <ResponsiveContainer width="100%" height={280}>
        <BarChart data={data} margin={{ top: 4, right: 8, left: -20, bottom: 0 }}>
          <defs>
            <linearGradient id="barGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="var(--primary)" stopOpacity={1} />
              <stop offset="100%" stopColor="var(--primary)" stopOpacity={0.4} />
            </linearGradient>
          </defs>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="var(--border)"
            vertical={false}
          />
          <XAxis
            dataKey="player"
            tick={{ fill: 'var(--muted-foreground)', fontSize: 12 }}
            axisLine={false}
            tickLine={false}
          />
          <YAxis
            tick={{ fill: 'var(--muted-foreground)', fontSize: 12 }}
            axisLine={false}
            tickLine={false}
          />
          <Tooltip
            content={<CustomTooltip />}
            cursor={{ fill: 'rgba(201,168,76,0.06)' }}
          />
          <Bar dataKey="value" fill="url(#barGradient)" radius={[4, 4, 0, 0]}>
            {data.map((entry, index) => (
              <Cell
                key={`cell-${entry.player}-${index}`}
                fill={index === 0 ? 'var(--primary)' : 'url(#barGradient)'}
                opacity={index === 0 ? 1 : 0.7}
              />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}