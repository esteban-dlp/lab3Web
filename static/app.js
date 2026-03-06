async function nextEpisode(id) {
  const res = await fetch(`/update?id=${id}`, { method: "POST" });
  if (!res.ok) {
    alert("No se pudo actualizar episodio.");
    return;
  }
  const data = await res.json();

  // Update numbers
  const cur = document.getElementById(`cur-${id}`);
  const tot = document.getElementById(`tot-${id}`);
  const pct = document.getElementById(`pct-${id}`);
  const bar = document.getElementById(`bar-${id}`);
  const done = document.getElementById(`done-${id}`);

  if (cur) cur.textContent = String(data.current);
  if (tot) tot.textContent = String(data.total);

  const percent = Math.round(data.progress);
  if (pct) pct.textContent = `${percent}%`;
  if (bar) bar.style.width = `${percent}%`;

  if (done) done.textContent = data.done ? "✅ Completada" : "";
}

async function deleteSeries(id) {
  const ok = confirm("¿Eliminar esta serie?");
  if (!ok) return;

  const res = await fetch(`/series?id=${id}`, { method: "DELETE" });
  if (!res.ok) {
    alert("No se pudo eliminar.");
    return;
  }

  const row = document.getElementById(`row-${id}`);
  if (row) row.remove();
}

async function saveRating(id) {
  const input = document.getElementById(`rate-${id}`);
  const rating = Number(input?.value ?? 0);

  if (!Number.isInteger(rating) || rating < 0 || rating > 10) {
    alert("Rating debe ser entero de 0 a 10.");
    return;
  }

  const res = await fetch(`/rating?id=${id}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ rating }),
  });

  if (!res.ok) {
    alert("No se pudo guardar rating.");
    return;
  }
}