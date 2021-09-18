include $(TOPDIR)/rules.mk

PKG_NAME:=hnsd
PKG_VERSION:=1
PKG_RELEASE:=1

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/handshake-org/hnsd.git
PKG_SOURCE_VERSION:=c60b2b49adf6de92f028670f1608493dbbad1578
PKG_MIRROR_HASH:=d2990a59bf091dda43a939ef76c5733be1fc1e727783c27b488fc49908e165ad

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)
PKG_BUILD_DEPENDS:=libunbound libuv

PKG_INSTALL:=1
PKG_BUILD_PARALLEL:=1

include $(INCLUDE_DIR)/package.mk
#include $(INCLUDE_DIR)/cmake.mk

define Package/hnsd
  SECTION:=net
  CATEGORY:=Network 
  TITLE:=hnsd
  URL:=https://hnsd.org/
  DEPENDS:=+libunbound
endef

define Package/hnsd/description
 hnsd
endef

#CMAKE_OPTIONS += -DBUILD_TESTING=OFF

define Build/Prepare
	$(Build/Prepare/Default)
	$(PKG_BUILD_DIR)/autogen.sh
endef

define Package/hnsd/install
	$(INSTALL_DIR) $(1)/usr/bin/
	$(INSTALL_BIN) $(PKG_INSTALL_DIR)/usr/bin/hnsd $(1)/usr/bin/
endef

$(eval $(call BuildPackage,hnsd))
