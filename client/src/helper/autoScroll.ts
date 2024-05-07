export function scrollToBottom(id: string) {
  const container = document.getElementById(id);
  console.log("should be scrolling", container);
  if (!container) {
    return;
  }
  console.log("ScrollTop before scrolling:", container.scrollTop);
  container.scrollTop = container.scrollHeight;
  console.log("ScrollTop after scrolling:", container.scrollTop);
}
