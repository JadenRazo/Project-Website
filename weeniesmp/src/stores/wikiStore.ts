import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface WikiPage {
  slug: string
  title: string
  content: string
  category: string
  lastUpdated: string
}

export interface WikiCategory {
  id: string
  name: string
  icon: string
  pages: WikiPage[]
}

// Static wiki content - imported at build time
const wikiContent: Record<string, { title: string; category: string; lastUpdated: string; content: string }> = {
  'getting-started': {
    title: 'Getting Started',
    category: 'basics',
    lastUpdated: '2026-01-21',
    content: `# Getting Started with WeenieSMP

Welcome to WeenieSMP! This guide will help you connect from **any platform** - Java Edition, Bedrock Edition (mobile/console/PC), or anywhere in between.

## Cross-Platform Play

WeenieSMP uses **Geyser + Floodgate** for true cross-platform play:
- ✅ Java and Bedrock players on the same server
- ✅ No Java Edition needed - use your Microsoft account!
- ✅ Works on PC, mobile, Xbox, PlayStation, and Switch

---

## Connection Instructions

### Java Edition (PC/Mac/Linux)

**Server Address:** \`play.weeniesmp.net\`

1. Minecraft → **Multiplayer** → **Add Server**
2. Enter address: \`play.weeniesmp.net\`
3. **Done** → Join!

---

### Bedrock Edition - Mobile & PC

**Server Address:** \`play.weeniesmp.net\` | **Port:** \`19011\`

#### Mobile (iOS/Android) & Windows 10/11

1. Minecraft → **Play** → **Servers** → **Add Server**
2. Enter server details:
   - **Name:** WeenieSMP
   - **Address:** \`play.weeniesmp.net\`
   - **Port:** \`19011\`
3. **Save** → Join!

> 💡 Bedrock usernames appear with a \`.\` prefix (e.g., \`.YourName\`)

---

### Console (PlayStation, Xbox, Nintendo Switch)

Consoles require extra steps due to manufacturer restrictions. Choose your platform:

#### PlayStation & Xbox

**Method: BedrockTogether App**

**Requirements:**
- Smartphone ([iOS](https://apps.apple.com/us/app/bedrocktogether/id1534593376) | [Android](https://play.google.com/store/apps/details?id=pl.extollite.bedrocktogetherapp)) on **same WiFi** as console
- [BedrockTogether](https://bedrocktogether.net/) app

**Steps:**
1. Download BedrockTogether on your phone
2. Connect phone & console to same WiFi
3. Open app, enter: \`play.weeniesmp.net:19011\`
4. Tap **"Start"** (keep app running!)
5. On console: **Play → Friends** → Find "BedrockTogether" in LAN Games
6. Join!

> ⚠️ **Keep your phone awake** - don't close the app or lock your screen while playing!

#### Nintendo Switch

**Method: DNS Redirect**

1. Switch → **Settings → Internet → DNS Settings → Manual**
2. Primary DNS: \`104.238.130.180\` | Secondary: \`8.8.8.8\`
3. Save & reconnect WiFi
4. Minecraft → **Featured Servers** → Join any server
5. BedrockConnect menu appears → Add server
6. Enter: \`play.weeniesmp.net:19011\`
7. Join!

> 💡 Revert DNS settings when not playing Minecraft

---

## First Steps

Once you're in:

1. **Read Rules** - Type \`/rules\` or visit [Rules](/rules)
2. **Find a Spot** - Use \`/rtp\` for random teleport
3. **Claim Land** - Use \`/claim\` to protect your base
4. **Set Home** - Use \`/sethome\` for quick returns
5. **Join Discord** - [discord.gg/weeniesmp](https://discord.gg/weeniesmp)

## Essential Commands

| Command | Description |
|---------|-------------|
| \`/spawn\` | Return to spawn |
| \`/rtp\` | Random teleport |
| \`/sethome\` / \`/home\` | Set/teleport to home |
| \`/claim\` | Claim your chunk |
| \`/balance\` | Check your money |

[Full Commands List](/wiki/commands)

---

## Quick Tips & Common Issues

### ✅ Do's
- ✅ **Double-check** server address and port
- ✅ **Update** Minecraft to latest version
- ✅ **Use same WiFi** for phone/console (BedrockTogether)
- ✅ **Keep app running** when using BedrockTogether
- ✅ **Ask for help** in Discord if stuck

### ❌ Don'ts
- ❌ **Don't** forget the port: \`19011\` for Bedrock Edition
- ❌ **Don't** close BedrockTogether app while playing
- ❌ **Don't** use guest/isolation WiFi networks
- ❌ **Don't** skip updating Minecraft to latest version

### Common Issues

**Can't connect?**
- Verify correct address: \`play.weeniesmp.net\` (Java) or \`play.weeniesmp.net:19011\` (Bedrock)
- Restart Minecraft
- Check internet connection

**Console not finding server?**
- Phone & console on same WiFi?
- BedrockTogether app still running?
- Try restarting the app

**Authentication error?**
- Check [Mojang Status](https://help.minecraft.net/hc/en-us/articles/360058525452-Minecraft-Services-Status)
- Wait a few minutes and retry

---

## Need More Help?

**Can't connect or having issues?**
- 💬 Ask in-game chat - our community is friendly!
- 🎫 **Open a ticket** in [Discord](https://discord.gg/weeniesmp)
- 📖 Check other wiki guides for specific features

**Learn More:**
- [Commands](/wiki/commands) - Full command list
- [Economy](/wiki/economy) - Earning money
- [Land Claims](/wiki/claims) - Protect your builds
- [Jobs](/wiki/jobs) - Work and earn

---

**Welcome to WeenieSMP!** 🎉

*Server powered by [Geyser](https://geysermc.org/) + [Floodgate](https://geysermc.org/wiki/floodgate/) for cross-platform play*
`
  },
  'commands': {
    title: 'Commands',
    category: 'basics',
    lastUpdated: '2025-01-15',
    content: `# Command Reference

Complete list of commands available on WeenieSMP.

## General Commands

### Navigation
| Command | Description | Permission |
|---------|-------------|------------|
| \`/spawn\` | Teleport to server spawn | Everyone |
| \`/sethome [name]\` | Set a home location | Everyone |
| \`/home [name]\` | Teleport to your home | Everyone |
| \`/delhome <name>\` | Delete a home | Everyone |
| \`/homes\` | List all your homes | Everyone |
| \`/rtp\` | Random teleport in the wild | Everyone |
| \`/tpa <player>\` | Request to teleport to a player | Everyone |
| \`/tpaccept\` | Accept a teleport request | Everyone |
| \`/tpdeny\` | Deny a teleport request | Everyone |
| \`/back\` | Return to your last location | Everyone |

### Communication
| Command | Description | Permission |
|---------|-------------|------------|
| \`/msg <player> <message>\` | Private message a player | Everyone |
| \`/r <message>\` | Reply to last message | Everyone |
| \`/mail send <player> <message>\` | Send offline mail | Everyone |
| \`/mail read\` | Read your mail | Everyone |

### Economy
| Command | Description | Permission |
|---------|-------------|------------|
| \`/balance\` or \`/bal\` | Check your balance | Everyone |
| \`/pay <player> <amount>\` | Send money to a player | Everyone |
| \`/baltop\` | View richest players | Everyone |

### Utility
| Command | Description | Permission |
|---------|-------------|------------|
| \`/kit\` | View available kits | Everyone |
| \`/kit <name>\` | Claim a kit | Everyone |
| \`/near\` | See nearby players | Everyone |
| \`/playtime\` | Check your playtime | Everyone |
| \`/seen <player>\` | Check when a player was last online | Everyone |

## Claiming Commands

See the [Claims Guide](/wiki/claims) for detailed information.

| Command | Description | Permission |
|---------|-------------|------------|
| \`/claim\` | Claim the chunk you're in | Everyone |
| \`/unclaim\` | Unclaim your chunk | Everyone |
| \`/claimlist\` | List all your claims | Everyone |
| \`/trust <player>\` | Trust a player in your claim | Everyone |
| \`/untrust <player>\` | Remove trust from a player | Everyone |
| \`/trustlist\` | View trusted players | Everyone |

## VIP Commands

These commands are available to players with VIP rank or higher.

| Command | Description | Rank |
|---------|-------------|------|
| \`/hat\` | Wear item as a hat | VIP+ |
| \`/craft\` | Open crafting table anywhere | VIP+ |
| \`/ec\` | Open ender chest anywhere | MVP+ |
| \`/nick\` | Set a nickname | PRO |
| \`/fly\` | Toggle creative flight | PRO |
`
  },
  'economy': {
    title: 'Economy',
    category: 'features',
    lastUpdated: '2025-01-15',
    content: `# Economy System

WeenieSMP features a balanced economy system designed to enhance gameplay without pay-to-win elements.

## Currency

The server uses **Weenies** ($) as the main currency.

## Earning Money

### Jobs
The primary way to earn money is through jobs. See the [Jobs Guide](/wiki/jobs) for more details.

### Selling Items
Sell items at the server shop:
- \`/shop\` - Open the shop GUI
- \`/sell hand\` - Sell the item in your hand
- \`/sell all\` - Sell all sellable items in your inventory

### Voting
Vote for the server to earn rewards:
- \`/vote\` - View voting links
- Each vote rewards you with money and vote keys

### Events
Participate in server events for cash prizes and rare items.

## Spending Money

### Claiming Land
Claiming chunks costs money. The cost increases exponentially with each claim:
- First claim: $100
- Second claim: $200
- And so on...

### Player Shops
Create your own shop or buy from other players:
- \`/chestshop\` - View chest shop help
- Place a sign on a chest to create a shop

### Server Shop
Buy useful items from the server:
- Spawn eggs
- Special blocks
- Cosmetic items

## Balance Commands

| Command | Description |
|---------|-------------|
| \`/balance\` | Check your balance |
| \`/pay <player> <amount>\` | Send money |
| \`/baltop\` | Leaderboard |

## Tips for Making Money

1. **Choose profitable jobs** - Mining and farming pay well
2. **Sell valuable items** - Diamonds, netherite, and rare items sell high
3. **Vote daily** - Consistent voting adds up
4. **Trade with players** - Sometimes player trades beat shop prices
5. **Complete quests** - Daily and weekly quests offer good rewards
`
  },
  'claims': {
    title: 'Land Claims',
    category: 'features',
    lastUpdated: '2025-01-15',
    content: `# Land Claiming System

Protect your builds from griefing with our chunk-based claiming system.

## How Claims Work

Claims are chunk-based (16x16 blocks, full height). When you claim a chunk:
- Only you can build, break, or open containers
- You can trust other players
- Mobs can still spawn naturally
- Your builds are protected even when offline

## Claiming Land

### Basic Claiming
1. Stand in the chunk you want to claim
2. Run \`/claim\`
3. Pay the claim cost

### Viewing Claims
- \`/claimlist\` - View all your claims
- Claims show particle borders when you enter them

### Unclaiming
- \`/unclaim\` - Remove claim on current chunk (50% refund)

## Claim Costs

Claims use exponential pricing to prevent land hoarding:

| Claim # | Cost |
|---------|------|
| 1 | $100 |
| 2 | $200 |
| 3 | $400 |
| 4 | $800 |
| 5 | $1,600 |
| ... | Doubles each time |

## Trust System

Share access to your claims with other players:

| Command | Description |
|---------|-------------|
| \`/trust <player>\` | Give full build access |
| \`/containertrust <player>\` | Chest access only |
| \`/accesstrust <player>\` | Door/button access only |
| \`/untrust <player>\` | Remove all trust |
| \`/trustlist\` | View trusted players |

### Trust Levels

1. **Access Trust** - Use doors, buttons, levers
2. **Container Trust** - Open chests, furnaces, etc.
3. **Build Trust** - Full building access

## Claim Limits

| Rank | Max Claims |
|------|------------|
| Default | 10 |
| VIP | 20 |
| MVP | 35 |
| PRO | 50 |

## Tips

- **Plan your base** - Claim strategically to save money
- **Trust wisely** - Only trust players you know
- **Claim borders** - Be aware of chunk boundaries when building
- **Check claims** - Use \`/claimlist\` to manage your territory
`
  },
  'jobs': {
    title: 'Jobs',
    category: 'features',
    lastUpdated: '2025-01-15',
    content: `# Jobs System

Earn money by doing activities you already enjoy in Minecraft.

## Available Jobs

### Miner
Earn money by mining ores and stone.

**Top Earners:**
- Diamond Ore: $15
- Ancient Debris: $50
- Emerald Ore: $12
- Gold Ore: $8
- Iron Ore: $5

### Woodcutter
Earn money by chopping trees.

**Pays for:**
- All log types
- Bonus for stripping logs
- Extra for jungle and dark oak (larger trees)

### Farmer
Earn money by harvesting crops.

**Top Earners:**
- Nether Wart: $3
- Wheat: $1
- Carrots: $1
- Potatoes: $1
- Melons/Pumpkins: $0.5 each

### Hunter
Earn money by killing mobs.

**Top Earners:**
- Wither: $500
- Ender Dragon: $1000
- Elder Guardian: $100
- Warden: $200
- Regular hostile mobs: $1-5

### Fisherman
Earn money by fishing.

**Pays for:**
- Fish catches
- Treasure items (bonus)
- Junk removal (small bonus)

### Builder
Earn money by placing blocks.

**Notes:**
- Only earns in claimed land
- Prevents exploit of break/place cycles
- Decorative blocks pay more

## Joining a Job

\`\`\`
/jobs browse    - View available jobs
/jobs join <job> - Join a job
/jobs leave <job> - Leave a job
/jobs stats     - View your job stats
\`\`\`

## Job Limits

| Rank | Max Jobs |
|------|----------|
| Default | 2 |
| VIP | 3 |
| MVP | 4 |
| PRO | 5 |

## Leveling Up

As you work, you gain job XP and level up:
- Higher levels = Better pay
- Level 10: 10% bonus
- Level 25: 25% bonus
- Level 50: 50% bonus
- Level 100: 100% bonus (2x pay!)

## Tips

1. **Specialize early** - Focus on 1-2 jobs to level faster
2. **Match your playstyle** - Pick jobs you naturally do
3. **Check prices** - Use \`/jobs info <job>\` to see payouts
4. **Level matters** - Stick with jobs to earn more over time
`
  }
}

// Build categories from wiki content
const categoryDefinitions: Record<string, { name: string; icon: string }> = {
  basics: { name: 'Getting Started', icon: 'book-open' },
  features: { name: 'Server Features', icon: 'star' }
}

export const useWikiStore = defineStore('wiki', () => {
  const currentPage = ref<WikiPage | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')

  const categories = computed<WikiCategory[]>(() => {
    const categoryMap: Record<string, WikiPage[]> = {}

    // Group pages by category
    for (const [slug, page] of Object.entries(wikiContent)) {
      const categoryId = page.category
      if (!categoryMap[categoryId]) {
        categoryMap[categoryId] = []
      }
      categoryMap[categoryId].push({
        slug,
        title: page.title,
        content: page.content,
        category: categoryId,
        lastUpdated: page.lastUpdated
      })
    }

    // Build category array
    return Object.entries(categoryMap).map(([id, pages]) => ({
      id,
      name: categoryDefinitions[id]?.name ?? id,
      icon: categoryDefinitions[id]?.icon ?? 'file-text',
      pages: pages.sort((a, b) => a.title.localeCompare(b.title))
    }))
  })

  const allPages = computed<WikiPage[]>(() => {
    return categories.value.flatMap(cat => cat.pages)
  })

  const searchResults = computed<WikiPage[]>(() => {
    if (!searchQuery.value.trim()) return []

    const query = searchQuery.value.toLowerCase()
    return allPages.value.filter(page =>
      page.title.toLowerCase().includes(query) ||
      page.content.toLowerCase().includes(query)
    ).slice(0, 10)
  })

  async function fetchCategories(): Promise<WikiCategory[]> {
    // Categories are computed from static content
    return categories.value
  }

  async function fetchPage(slug: string): Promise<WikiPage | null> {
    loading.value = true
    error.value = null

    try {
      // Simulate async for consistency with API pattern
      await new Promise(resolve => setTimeout(resolve, 50))

      const pageData = wikiContent[slug]
      if (!pageData) {
        error.value = 'Page not found'
        currentPage.value = null
        return null
      }

      const page: WikiPage = {
        slug,
        title: pageData.title,
        content: pageData.content,
        category: pageData.category,
        lastUpdated: pageData.lastUpdated
      }

      currentPage.value = page
      return page
    } catch (e) {
      error.value = 'Failed to load page'
      currentPage.value = null
      return null
    } finally {
      loading.value = false
    }
  }

  function getAdjacentPages(slug: string): { prev: WikiPage | null; next: WikiPage | null } {
    const pages = allPages.value
    const currentIndex = pages.findIndex(p => p.slug === slug)

    return {
      prev: currentIndex > 0 ? pages[currentIndex - 1] : null,
      next: currentIndex < pages.length - 1 ? pages[currentIndex + 1] : null
    }
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  return {
    currentPage,
    loading,
    error,
    categories,
    allPages,
    searchQuery,
    searchResults,
    fetchCategories,
    fetchPage,
    getAdjacentPages,
    setSearchQuery
  }
})
