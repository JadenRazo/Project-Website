import { useState } from 'react'
import { Send, Mail, Github, Linkedin, CheckCircle, Loader2 } from 'lucide-react'
import { api } from '../../../utils/apiConfig'

const socialLinks = [
  { icon: Github, href: 'https://github.com/JadenRazo', label: 'GitHub' },
  { icon: Linkedin, href: 'https://jadenrazo.dev/s/linkedin', label: 'LinkedIn' },
  { icon: Mail, href: 'mailto:contact@jadenrazo.dev', label: 'Email' },
]

export default function Contact() {
  const [isSubmitted, setIsSubmitted] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: '',
    message: '',
    website: '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      await api.post('/api/v1/contact', {
        name: formData.name,
        email: formData.email,
        subject: formData.subject,
        message: formData.message,
        website: formData.website,
      }, { skipAuth: true })

      setIsSubmitted(true)
      setFormData({ name: '', email: '', subject: '', message: '', website: '' })
      setTimeout(() => setIsSubmitted(false), 5000)
    } catch (err) {
      setError(
        err instanceof Error && err.message !== 'Request failed'
          ? err.message
          : 'Failed to send message. Please try again later.'
      )
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }))
    if (error) setError(null)
  }

  return (
    <section id="contact" className="relative w-full py-12 sm:py-16 md:py-20 lg:py-28">
      <div className="portfolio-container flex flex-col items-center">
        <div
          className="max-w-3xl text-center mb-6 md:mb-8 lg:mb-10"
        >
          <p className="mb-3 font-mono text-xs uppercase tracking-[0.2em] text-primary">
            Start a conversation
          </p>
          <h2 className="font-display text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold tracking-[-0.04em] mb-3 lg:mb-5 text-text-primary">
            Let&apos;s build something dependable.
          </h2>
          <p className="text-[15px] leading-7 text-text-secondary sm:text-base lg:text-lg">
            Hiring for a cloud, DevOps, platform, or SRE team? Tell me about the
            system, the reliability problem, and what ownership looks like.
          </p>
          <a href="mailto:contact@jadenrazo.dev" className="mt-3 inline-block text-sm font-medium text-primary hover:underline">
            contact@jadenrazo.dev
          </a>
        </div>

        <div
          className="border border-border bg-background-secondary p-4 sm:p-6 md:p-7 lg:p-8 w-full max-w-xl md:max-w-lg lg:max-w-2xl"
        >
          {isSubmitted ? (
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <CheckCircle className="w-10 h-10 sm:w-12 sm:h-12 text-primary mb-3" />
              <h3 className="text-lg sm:text-xl font-bold mb-1 text-text-primary">Message sent</h3>
              <p className="text-[15px] leading-6 text-text-secondary">
                Thanks for reaching out. I'll get back to you soon.
              </p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-3 sm:space-y-4 lg:space-y-5">
              {error && (
                <div className="px-3 py-2 bg-red-500/10 border border-red-500/30 rounded-lg text-red-400 text-sm">
                  {error}
                </div>
              )}
              <div aria-hidden="true" style={{ position: 'absolute', left: '-9999px', opacity: 0, height: 0, overflow: 'hidden' }}>
                <input
                  type="text"
                  name="website"
                  value={formData.website}
                  onChange={handleChange}
                  tabIndex={-1}
                  autoComplete="off"
                />
              </div>
              <div className="grid sm:grid-cols-2 gap-3 sm:gap-4">
                <div>
                  <label className="sr-only" htmlFor="contact-name">Your name</label>
                  <input
                    type="text"
                    id="contact-name"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    disabled={isLoading}
                    className="w-full px-4 py-3 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors text-base sm:text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                    placeholder="Your name"
                  />
                </div>
                <div>
                  <label className="sr-only" htmlFor="contact-email">Your email</label>
                  <input
                    type="email"
                    id="contact-email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    disabled={isLoading}
                    className="w-full px-4 py-3 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors text-base sm:text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                    placeholder="Your email"
                  />
                </div>
              </div>
              <div>
                <label className="sr-only" htmlFor="contact-subject">Subject</label>
                <input
                  type="text"
                  id="contact-subject"
                  name="subject"
                  value={formData.subject}
                  onChange={handleChange}
                  disabled={isLoading}
                  className="w-full px-4 py-3 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors text-base sm:text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                  placeholder="Role or team (optional)"
                />
              </div>
              <div>
                <label className="sr-only" htmlFor="contact-message">Your message</label>
                <textarea
                  id="contact-message"
                  name="message"
                  value={formData.message}
                  onChange={handleChange}
                  required
                  disabled={isLoading}
                  rows={3}
                  className="w-full px-4 py-3 bg-surface border border-border rounded-lg focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors resize-none text-base sm:text-sm text-text-primary placeholder:text-text-muted disabled:opacity-50"
                  placeholder="What are you building or hiring for?"
                />
              </div>
              <button type="submit" className="btn-primary w-full" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <Loader2 size={16} className="animate-spin" />
                    <span>Sending...</span>
                  </>
                ) : (
                  <>
                    <Send size={16} />
                    <span>Send Message</span>
                  </>
                )}
              </button>
            </form>
          )}
        </div>

        <div
          className="flex items-center justify-center gap-4 mt-5"
        >
          {socialLinks.map((social) => (
            <a
              key={social.label}
              href={social.href}
              target="_blank"
              rel="noopener noreferrer"
              className="p-3 min-w-[44px] min-h-[44px] flex items-center justify-center border border-border hover:border-primary active:bg-surface-hover transition-colors group"
              aria-label={social.label}
            >
              <social.icon className="w-5 h-5 text-text-secondary group-hover:text-primary transition-colors" />
            </a>
          ))}
        </div>

        <p className="mt-5 text-center text-[15px] leading-6 text-text-muted">
          &copy; {new Date().getFullYear()} Jaden Razo. All rights reserved.
        </p>
      </div>
    </section>
  )
}
