// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  // Enable experimental app directory if needed (Next 14 default)
  experimental: {
    appDir: true,
  },
};
module.exports = nextConfig;
