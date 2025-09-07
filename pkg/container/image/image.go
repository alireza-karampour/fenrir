package image

import (
	"context"

	"github.com/containers/buildah"
	"github.com/containers/image/docker/reference"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/archive"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

type Base struct {
	Ctx           context.Context
	Store         storage.Store
	From          string
	Builder       *buildah.Builder
	ImageID       *string
	ConanicalName reference.Canonical
	Digest        *digest.Digest
	Logger        *logrus.Logger
}

func From(ctx context.Context, image string) (*Base, error) {
	storeOpt, err := storage.DefaultStoreOptions()
	if err != nil {
		panic(err)
	}
	store, err := storage.GetStore(storeOpt)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	builder, err := buildah.NewBuilder(ctx, store, buildah.BuilderOptions{
		FromImage: image,
		Registry:  "docker.io",
		Logger:    logger,
	})
	if err != nil {
		return nil, err
	}
	return &Base{
		Ctx:     ctx,
		From:    image,
		Builder: builder,
		Store:   store,
	}, nil
}

func (b *Base) Copy(contextDir string, src string, dst string) (*Base, error) {
	err := b.Builder.Add(dst, false, buildah.AddAndCopyOptions{
		ContextDir: contextDir,
		Link:       true,
	}, src)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Base) WorkDir(path string) *Base {
	b.Builder.SetWorkDir(path)
	return b
}

func (b *Base) Run(cmdline ...string) (*Base, error) {
	err := b.Builder.Run(cmdline, buildah.RunOptions{
		Logger: b.Logger,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Base) Env(key string, value string) *Base {
	b.Builder.SetEnv(key, value)
	return b
}

func (b *Base) Entrypoint(entry ...string) *Base {
	b.Builder.SetEntrypoint(entry)
	return b
}

func (b *Base) Cmd(cmdline ...string) *Base {
	b.Builder.SetCmd(cmdline)
	return b
}

func (b *Base) Build(tag string) (*Base, error) {
	imageRef, err := is.Transport.ParseStoreReference(b.Store, tag)
	if err != nil {
		return nil, err
	}

	imageId, conanical, digest, err := b.Builder.Commit(b.Ctx, imageRef, buildah.CommitOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return nil, err
	}

	b.ImageID = &imageId
	b.ConanicalName = conanical
	b.Digest = &digest

	return b, nil
}
