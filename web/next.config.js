/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'picsum.photos',
      },
      {
        protocol: 'https',
        hostname: 'next-bazaar.s3.ap-northeast-1.amazonaws.com',
      },
    ],
  },
}

module.exports = nextConfig
