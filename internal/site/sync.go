package site

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
)

// Sync copies every file listed in the manifest from the wiki into docsDir,
// keeping the wiki's relative paths (so relative links between guides keep
// working). Only files whose content changed are written. It returns the
// paths (relative to docsDir) that were created or updated.
func Sync(m *Manifest, docsDir string) ([]string, error) {
	var changed []string
	for _, rel := range m.Files() {
		src := filepath.Join(m.WikiRoot(), rel)
		dst := filepath.Join(docsDir, rel)
		in, err := os.ReadFile(src)
		if err != nil {
			return changed, err
		}
		if cur, err := os.ReadFile(dst); err == nil && bytes.Equal(cur, in) {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return changed, err
		}
		// Write to a temporary file and rename, so a watcher never sees half a file.
		tmp := dst + ".tmp"
		if err := os.WriteFile(tmp, in, 0o644); err != nil {
			return changed, err
		}
		if err := os.Rename(tmp, dst); err != nil {
			return changed, err
		}
		changed = append(changed, rel)
	}
	return changed, nil
}

// fingerprint is a cheap change detector: modification time + size per file.
type fingerprint map[string][2]int64

func takeFingerprint(root string, rels []string) fingerprint {
	fp := fingerprint{}
	for _, rel := range rels {
		if st, err := os.Stat(filepath.Join(root, rel)); err == nil {
			fp[rel] = [2]int64{st.ModTime().UnixNano(), st.Size()}
		}
	}
	return fp
}

// treeFingerprint covers every file under dir (for templates and assets).
func treeFingerprint(dir string) fingerprint {
	fp := fingerprint{}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if st, err := d.Info(); err == nil {
				fp[p] = [2]int64{st.ModTime().UnixNano(), st.Size()}
			}
		}
		return nil
	})
	return fp
}

// diff returns the keys whose fingerprint changed, appeared, or disappeared.
func (a fingerprint) diff(b fingerprint) []string {
	var out []string
	for k, v := range b {
		if a[k] != v {
			out = append(out, k)
		}
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			out = append(out, k)
		}
	}
	return out
}
