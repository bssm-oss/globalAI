async function bootstrap() {
  const summary = document.querySelector('#summary');
  const sourceList = document.querySelector('#source-list');
  const sourceTitle = document.querySelector('#source-title');
  const sourceMeta = document.querySelector('#source-meta');
  const sourceContent = document.querySelector('#source-content');

  const response = await fetch('/api/sources');
  const payload = await response.json();
  const sources = payload.sources || [];
  const generatedAt = payload.generatedAt ? new Date(payload.generatedAt).toLocaleString() : 'unknown';

  const total = document.createElement('strong');
  total.textContent = `${payload.totalSources || sources.length} sources`;
  const root = document.createElement('p');
  root.textContent = `root · ${payload.root || 'unknown'}`;
  const generated = document.createElement('p');
  generated.textContent = `updated · ${generatedAt}`;
  summary.replaceChildren(total, root, generated);

  if (sources.length === 0) {
    sourceTitle.textContent = 'No sources found';
    sourceMeta.textContent = 'Add an allowlisted file such as AGENTS.md, CLAUDE.md, .claude/, .cursor/rules/, or .sisyphus/.';
    sourceContent.textContent = '';
    return;
  }

  function render(selectedSource) {
    sourceTitle.textContent = selectedSource.label;
    sourceMeta.textContent = `${selectedSource.family} · ${selectedSource.relativePath} · ${selectedSource.sizeBytes} bytes`;
    sourceContent.textContent = selectedSource.content;
  }

  sources.forEach((item, index) => {
    const listItem = document.createElement('li');
    const button = document.createElement('button');
    const label = document.createElement('strong');
    const meta = document.createElement('small');
    button.type = 'button';
    button.className = 'source-button';
    label.textContent = item.label;
    meta.textContent = `${item.family} · ${item.relativePath}`;
    button.append(label, meta);
    button.addEventListener('click', () => {
      document.querySelectorAll('.source-button').forEach((entry) => entry.classList.remove('is-active'));
      button.classList.add('is-active');
      render(item);
    });

    if (index === 0) {
      button.classList.add('is-active');
      render(item);
    }

    listItem.appendChild(button);
    sourceList.appendChild(listItem);
  });
}

bootstrap().catch((error) => {
  const sourceTitle = document.querySelector('#source-title');
  const sourceMeta = document.querySelector('#source-meta');
  sourceTitle.textContent = 'Load failed';
  sourceMeta.textContent = error.message;
});
