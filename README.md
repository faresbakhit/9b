# 9B – A link aggregator you can self-host

<img align="left" src="art/icon-with-tech-mascots.svg" width="256px">

9B (Nine B) is a [free] link aggregator/social news aggregation/discussion forum
service that you can self-host easily for any purpose you want (e.g. in-house
forum for your team, public forum for a niche topic, etc).

It is similar to websites you've probably encountered before like [Reddit],
[Lemmy], and [Hacker News], however 9B differs from those that you can:

- Build it and self-host it yourself in a matter of seconds.
- Have permissions and roles for each user.
- Customize the interface for your liking.

[free]: https://www.gnu.org/philosophy/free-sw.en.html
[Reddit]: https://old.reddit.com
[Lemmy]: https://join-lemmy.org/
[Hacker News]: https://news.ycombinator.com

## Routes

- `GET /signup`
- `POST /signup`
- `GET /login`
- `POST /login`
- `GET /u/{name}`
- `POST /u/{name}`
- `GET /protected`
