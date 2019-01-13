
# About

tcp_replay is copycat of [goreplay](https://github.com/buger/goreplay), works on tcp tracffic.

It's currently a toy project and not tested as well .

All credit goes to Leonid Bugaev, [@buger](https://twitter.com/buger), https://leonsbox.com



# Usage

```
# Running as non root user
# Test
sudo ./tcp_replay --input-tcp 127.0.0.1:4000 --output-stdout
