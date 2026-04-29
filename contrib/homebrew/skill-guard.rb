class SkillGuard < Formula
  desc "Security scanner for AI skill files"
  homepage "https://github.com/wjames2000/skill-guard"
  version "0.2.0"
  license "MIT"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/wjames2000/skill-guard/releases/download/v#{version}/skill-guard-darwin-arm64"
      sha256 "PLACEHOLDER_RUN: shasum -a 256 dist/skill-guard-darwin-arm64"
    else
      url "https://github.com/wjames2000/skill-guard/releases/download/v#{version}/skill-guard-darwin-amd64"
      sha256 "PLACEHOLDER_RUN: shasum -a 256 dist/skill-guard-darwin-amd64"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/wjames2000/skill-guard/releases/download/v#{version}/skill-guard-linux-arm64"
      sha256 "PLACEHOLDER_RUN: shasum -a 256 dist/skill-guard-linux-arm64"
    else
      url "https://github.com/wjames2000/skill-guard/releases/download/v#{version}/skill-guard-linux-amd64"
      sha256 "PLACEHOLDER_RUN: shasum -a 256 dist/skill-guard-linux-amd64"
    end
  end

  def install
    bin.install Dir["skill-guard-*"].first => "skill-guard"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/skill-guard --version")
  end
end
