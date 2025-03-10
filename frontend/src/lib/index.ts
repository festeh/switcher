

export function parseRun(run: string): { command: string; arg: string } {
  const parts = run.split(" ");
  const command = parts[0];
  const arg = parts[1];
  return { command, arg };
}
