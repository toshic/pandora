package guns

import (
	"github.com/spf13/afero"
	"github.com/yandex/pandora/components/guns/dummy"
	"github.com/yandex/pandora/core/register"
)

func Import(fs afero.Fs) {
	register.Gun("dummy", dummy.NewGun, dummy.DefaultGunConfig)
}
