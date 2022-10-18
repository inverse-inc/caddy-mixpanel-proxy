# Mixpanel plugin for PacketFence

This is a plugin used by the PacketFence team to gather anonymous analytics data and forward it to Mixpanel.

# Usage

Example Caddyfile

```
{
  order mixpanel_proxy before reverse_proxy
}

analytics.zammitcorp.com {
  mixpanel_proxy

  reverse_proxy https://api.mixpanel.com
}
```
