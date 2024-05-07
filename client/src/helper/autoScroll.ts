export function scrollToBottom(id: string) {
  const container = document.getElementById(id);
  if (!container) {
    return;
  }
  container.scrollTop = container.scrollHeight;
}
