include $(TOPDIR)/rules.mk

PKG_NAME:=dvpn-interface
PKG_VERSION:=1.0.0
PKG_RELEASE:=1

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/audi70r/dvpn-openwrt.git
PKG_SOURCE_VERSION:=cf0373f8b1b376594478623c49a5dd7e278b7dff
PKG_MIRROR_HASH:=943b98fe413494f0f730c6defdff518391f4685ad838f65579c62d1ac831c5fd

PKG_MAINTAINER:=your_name

PKG_BUILD_DIR:=$(BUILD_DIR)/dvpn-interface-$(PKG_VERSION)
PKG_BUILD_DEPENDS:=golang/host dvpn-node
PKG_BUILD_PARALLEL:=1
PKG_USE_MIPS16:=0

GO_PKG:=github.com/audi70r/dvpn-openwrt
GO_PKG_BUILD_PKG:=github.com/audi70r/dvpn-openwrt
GO_PKG_LDFLAGS:=-s -w

include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/nls.mk
include ../../lang/golang/golang-package.mk

define Package/dvpn-interface
	SECTION:=net
	CATEGORY:=Network
	DEPENDS:=$(GO_ARCH_DEPENDS)
	TITLE:=dvpn-interface
	URL:=https://github.com/audi70r/dvpn-openwrt
endef

define Package/dvpn-interface/install
	$(call GoPackage/Package/Install/Bin,$(1))
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./files/dvpn-interface.init $(1)/etc/init.d/dvpn-interface
endef

$(eval $(call GoBinPackage,dvpn-interface))
$(eval $(call BuildPackage,dvpn-interface))
